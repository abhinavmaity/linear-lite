package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type IssueRepositoryDB struct {
	db *gorm.DB
}

func NewIssueRepository(db *gorm.DB) *IssueRepositoryDB {
	return &IssueRepositoryDB{db: db}
}

func (r *IssueRepositoryDB) List(ctx context.Context, filter IssueListFilter) ([]models.Issue, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Issue{})
	query = applyIssueListFilters(query, filter)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []models.Issue{}, 0, nil
	}

	sortBy := filter.By
	if sortBy == "" {
		sortBy = "updated_at"
	}
	order := strings.ToLower(strings.TrimSpace(filter.Order))
	if order != "asc" {
		order = "desc"
	}

	orderClause := "updated_at desc"
	switch sortBy {
	case "identifier":
		orderClause = "identifier " + order
	case "title":
		orderClause = "title " + order
	case "status":
		orderClause = "status " + order
	case "priority":
		orderClause = "priority " + order
	case "created_at":
		orderClause = "created_at " + order
	case "updated_at":
		orderClause = "updated_at " + order
	}

	var issues []models.Issue
	err := query.Order(orderClause).
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Find(&issues).Error
	if err != nil {
		return nil, 0, err
	}
	return issues, total, nil
}

func (r *IssueRepositoryDB) FindByID(ctx context.Context, id string, includeArchived bool) (*models.Issue, error) {
	query := r.db.WithContext(ctx).Model(&models.Issue{}).Where("id = ?", strings.TrimSpace(id))
	if !includeArchived {
		query = query.Where("archived_at IS NULL")
	}

	var issue models.Issue
	err := query.First(&issue).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

type CreateIssueInput struct {
	Title       string
	Description *string
	Status      string
	Priority    string
	ProjectID   string
	SprintID    *string
	AssigneeID  *string
	CreatedBy   string
	LabelIDs    []string
}

func (r *IssueRepositoryDB) CreateWithRelations(ctx context.Context, input CreateIssueInput) (*models.Issue, error) {
	var created *models.Issue
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var projectLock struct {
			ID              string
			Key             string
			NextIssueNumber int
		}
		if err := tx.Raw(
			`SELECT id, key, next_issue_number FROM projects WHERE id = ? FOR UPDATE`,
			input.ProjectID,
		).Scan(&projectLock).Error; err != nil {
			return err
		}
		if projectLock.ID == "" {
			return ErrNotFound
		}

		identifier := fmt.Sprintf("%s-%d", projectLock.Key, projectLock.NextIssueNumber)
		now := time.Now().UTC()
		issue := &models.Issue{
			ID:         uuid.NewString(),
			Identifier: identifier,
			Title:      input.Title,
			Status:     input.Status,
			Priority:   input.Priority,
			ProjectID:  input.ProjectID,
			CreatedBy:  input.CreatedBy,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		issue.Description = input.Description
		issue.SprintID = input.SprintID
		issue.AssigneeID = input.AssigneeID

		if err := tx.Create(issue).Error; err != nil {
			if isUniqueViolationConstraint(err, "uq_issues_identifier") {
				return ErrConflict
			}
			return err
		}

		if len(input.LabelIDs) > 0 {
			links := make([]models.IssueLabel, 0, len(input.LabelIDs))
			for _, labelID := range input.LabelIDs {
				links = append(links, models.IssueLabel{
					IssueID:   issue.ID,
					LabelID:   labelID,
					CreatedAt: now,
				})
			}
			if err := tx.Create(&links).Error; err != nil {
				return err
			}
		}

		activities := []models.IssueActivity{
			{
				ID:        uuid.NewString(),
				IssueID:   issue.ID,
				UserID:    input.CreatedBy,
				Action:    models.IssueActivityCreated,
				CreatedAt: now,
			},
		}
		for _, labelID := range input.LabelIDs {
			newValue := labelID
			activities = append(activities, models.IssueActivity{
				ID:        uuid.NewString(),
				IssueID:   issue.ID,
				UserID:    input.CreatedBy,
				Action:    models.IssueActivityLabelAdded,
				FieldName: strPtr("labels"),
				NewValue:  &newValue,
				CreatedAt: now,
			})
		}
		if err := tx.Create(&activities).Error; err != nil {
			return err
		}

		if err := tx.Exec(
			`UPDATE projects SET next_issue_number = next_issue_number + 1 WHERE id = ?`,
			input.ProjectID,
		).Error; err != nil {
			return err
		}

		created = issue
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

type UpdateIssueInput struct {
	ID          string
	Title       *string
	Description **string
	Status      *string
	Priority    *string
	ProjectID   *string
	SprintID    **string
	AssigneeID  **string
	ActorID     string
	LabelIDs    *[]string
	Restore     bool
}

func (r *IssueRepositoryDB) UpdateWithRelations(ctx context.Context, input UpdateIssueInput) (*models.Issue, bool, error) {
	var updated *models.Issue
	var changed bool

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var issue models.Issue
		if err := tx.Where("id = ?", input.ID).First(&issue).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		now := time.Now().UTC()
		activities := make([]models.IssueActivity, 0)

		addFieldActivity := func(action, field string, oldValue, newValue *string) {
			activities = append(activities, models.IssueActivity{
				ID:        uuid.NewString(),
				IssueID:   issue.ID,
				UserID:    input.ActorID,
				Action:    action,
				FieldName: strPtr(field),
				OldValue:  oldValue,
				NewValue:  newValue,
				CreatedAt: now,
			})
		}

		if input.Restore {
			if issue.ArchivedAt == nil {
				return ErrConflict
			}
			issue.ArchivedAt = nil
			issue.ArchivedBy = nil
			activities = append(activities, models.IssueActivity{
				ID:        uuid.NewString(),
				IssueID:   issue.ID,
				UserID:    input.ActorID,
				Action:    models.IssueActivityRestored,
				CreatedAt: now,
			})
			changed = true
		}

		if input.Title != nil && strings.TrimSpace(*input.Title) != issue.Title {
			oldValue := issue.Title
			newValue := strings.TrimSpace(*input.Title)
			issue.Title = newValue
			addFieldActivity(models.IssueActivityTitleChanged, "title", &oldValue, &newValue)
			changed = true
		}
		if input.Description != nil {
			oldVal := issue.Description
			newVal := *input.Description
			if !strPtrEqual(oldVal, newVal) {
				issue.Description = newVal
				addFieldActivity(models.IssueActivityDescriptionChanged, "description", oldVal, newVal)
				changed = true
			}
		}
		if input.Status != nil && *input.Status != issue.Status {
			oldValue := issue.Status
			newValue := *input.Status
			issue.Status = newValue
			addFieldActivity(models.IssueActivityStatusChanged, "status", &oldValue, &newValue)
			changed = true
		}
		if input.Priority != nil && *input.Priority != issue.Priority {
			oldValue := issue.Priority
			newValue := *input.Priority
			issue.Priority = newValue
			addFieldActivity(models.IssueActivityPriorityChanged, "priority", &oldValue, &newValue)
			changed = true
		}
		if input.ProjectID != nil && *input.ProjectID != issue.ProjectID {
			oldValue := issue.ProjectID
			newValue := *input.ProjectID
			issue.ProjectID = newValue
			addFieldActivity(models.IssueActivityProjectChanged, "project_id", &oldValue, &newValue)
			changed = true
		}
		if input.SprintID != nil {
			oldVal := issue.SprintID
			newVal := *input.SprintID
			if !strPtrEqual(oldVal, newVal) {
				issue.SprintID = newVal
				addFieldActivity(models.IssueActivitySprintChanged, "sprint_id", oldVal, newVal)
				changed = true
			}
		}
		if input.AssigneeID != nil {
			oldVal := issue.AssigneeID
			newVal := *input.AssigneeID
			if !strPtrEqual(oldVal, newVal) {
				issue.AssigneeID = newVal
				addFieldActivity(models.IssueActivityAssigneeChanged, "assignee_id", oldVal, newVal)
				changed = true
			}
		}

		if changed {
			if err := tx.Save(&issue).Error; err != nil {
				return err
			}
		}

		if input.LabelIDs != nil {
			current, err := r.listLabelIDsByIssueIDTx(tx, issue.ID)
			if err != nil {
				return err
			}
			currentSet := make(map[string]struct{}, len(current))
			for _, id := range current {
				currentSet[id] = struct{}{}
			}
			nextSet := make(map[string]struct{}, len(*input.LabelIDs))
			for _, id := range *input.LabelIDs {
				nextSet[id] = struct{}{}
			}

			for id := range nextSet {
				if _, ok := currentSet[id]; !ok {
					val := id
					activities = append(activities, models.IssueActivity{
						ID:        uuid.NewString(),
						IssueID:   issue.ID,
						UserID:    input.ActorID,
						Action:    models.IssueActivityLabelAdded,
						FieldName: strPtr("labels"),
						NewValue:  &val,
						CreatedAt: now,
					})
					changed = true
				}
			}
			for id := range currentSet {
				if _, ok := nextSet[id]; !ok {
					val := id
					activities = append(activities, models.IssueActivity{
						ID:        uuid.NewString(),
						IssueID:   issue.ID,
						UserID:    input.ActorID,
						Action:    models.IssueActivityLabelRemoved,
						FieldName: strPtr("labels"),
						OldValue:  &val,
						CreatedAt: now,
					})
					changed = true
				}
			}

			if err := tx.Where("issue_id = ?", issue.ID).Delete(&models.IssueLabel{}).Error; err != nil {
				return err
			}
			if len(*input.LabelIDs) > 0 {
				links := make([]models.IssueLabel, 0, len(*input.LabelIDs))
				for _, labelID := range *input.LabelIDs {
					links = append(links, models.IssueLabel{
						IssueID:   issue.ID,
						LabelID:   labelID,
						CreatedAt: now,
					})
				}
				if err := tx.Create(&links).Error; err != nil {
					return err
				}
			}
		}

		if len(activities) > 0 {
			if err := tx.Create(&activities).Error; err != nil {
				return err
			}
		}

		updated = &issue
		return nil
	})
	if err != nil {
		return nil, false, err
	}
	return updated, changed, nil
}

func (r *IssueRepositoryDB) Archive(ctx context.Context, id string, actorID string) (bool, error) {
	var mutated bool
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var issue models.Issue
		if err := tx.Where("id = ?", id).First(&issue).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if issue.ArchivedAt != nil {
			mutated = false
			return nil
		}

		now := time.Now().UTC()
		issue.ArchivedAt = &now
		issue.ArchivedBy = &actorID
		if err := tx.Save(&issue).Error; err != nil {
			return err
		}
		activity := models.IssueActivity{
			ID:        uuid.NewString(),
			IssueID:   issue.ID,
			UserID:    actorID,
			Action:    models.IssueActivityArchived,
			CreatedAt: now,
		}
		if err := tx.Create(&activity).Error; err != nil {
			return err
		}
		mutated = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return mutated, nil
}

func (r *IssueRepositoryDB) ListLabelsByIssueIDs(ctx context.Context, issueIDs []string) (map[string][]models.Label, error) {
	out := make(map[string][]models.Label, len(issueIDs))
	if len(issueIDs) == 0 {
		return out, nil
	}

	type row struct {
		IssueID string
		models.Label
	}
	var rows []row
	err := r.db.WithContext(ctx).Table("issue_labels il").
		Select("il.issue_id, l.id, l.name, l.color, l.description, l.created_at, l.updated_at").
		Joins("JOIN labels l ON l.id = il.label_id").
		Where("il.issue_id IN ?", issueIDs).
		Order("l.name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.IssueID] = append(out[r.IssueID], r.Label)
	}
	return out, nil
}

func (r *IssueRepositoryDB) ListActivitiesByIssueID(ctx context.Context, issueID string) ([]models.IssueActivity, error) {
	var rows []models.IssueActivity
	err := r.db.WithContext(ctx).Model(&models.IssueActivity{}).
		Where("issue_id = ?", issueID).
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *IssueRepositoryDB) DashboardMetrics(ctx context.Context, userID string, doneSince time.Time) (DashboardMetrics, error) {
	var metrics DashboardMetrics
	// Raw SQL is used here to keep dashboard aggregation explicit and efficient.
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) FILTER (WHERE archived_at IS NULL)::int AS total_issues,
			COUNT(*) FILTER (WHERE archived_at IS NULL AND assignee_id = @user_id)::int AS my_issues,
			COUNT(*) FILTER (WHERE archived_at IS NULL AND status = 'in_progress')::int AS in_progress,
			COUNT(*) FILTER (WHERE archived_at IS NULL AND status = 'done' AND updated_at >= @done_since)::int AS done_this_week
		FROM issues
	`, map[string]any{
		"user_id":    strings.TrimSpace(userID),
		"done_since": doneSince.UTC(),
	}).Scan(&metrics).Error
	if err != nil {
		return DashboardMetrics{}, err
	}
	return metrics, nil
}

func (r *IssueRepositoryDB) DashboardActiveSprintID(ctx context.Context) (*string, error) {
	var sprintID string
	// Raw SQL gives deterministic active-sprint selection with explicit ordering.
	err := r.db.WithContext(ctx).Raw(`
		SELECT id
		FROM sprints
		WHERE status = 'active'
		ORDER BY end_date ASC, created_at ASC
		LIMIT 1
	`).Scan(&sprintID).Error
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(sprintID) == "" {
		return nil, nil
	}
	return &sprintID, nil
}

func (r *IssueRepositoryDB) ListRecentActivitiesForDashboard(ctx context.Context, limit int) ([]models.IssueActivity, error) {
	if limit <= 0 {
		return []models.IssueActivity{}, nil
	}
	var rows []models.IssueActivity
	// Raw SQL keeps the archive filter and ordering clear for dashboard activity.
	err := r.db.WithContext(ctx).Raw(`
		SELECT ia.id, ia.issue_id, ia.user_id, ia.action, ia.field_name, ia.old_value, ia.new_value, ia.created_at
		FROM issue_activities ia
		INNER JOIN issues i ON i.id = ia.issue_id
		WHERE i.archived_at IS NULL
		ORDER BY ia.created_at DESC
		LIMIT ?
	`, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *IssueRepositoryDB) listLabelIDsByIssueIDTx(tx *gorm.DB, issueID string) ([]string, error) {
	var ids []string
	err := tx.Model(&models.IssueLabel{}).
		Where("issue_id = ?", issueID).
		Pluck("label_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func applyIssueListFilters(query *gorm.DB, filter IssueListFilter) *gorm.DB {
	if !filter.IncludeArchived {
		query = query.Where("archived_at IS NULL")
	}
	if filter.Search != "" {
		search := strings.TrimSpace(filter.Search)
		pattern := "%" + search + "%"
		query = query.Where(
			"(identifier ILIKE ? OR search_vector @@ plainto_tsquery('english', ?))",
			pattern,
			search,
		)
	}
	if len(filter.Statuses) > 0 {
		query = query.Where("status IN ?", filter.Statuses)
	}
	if len(filter.Priorities) > 0 {
		query = query.Where("priority IN ?", filter.Priorities)
	}
	if filter.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *filter.AssigneeID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.SprintID != nil {
		query = query.Where("sprint_id = ?", *filter.SprintID)
	}
	if len(filter.LabelIDs) > 0 {
		mode := filter.LabelMode
		if mode == "" {
			mode = "any"
		}
		if mode == "all" {
			sub := query.Session(&gorm.Session{}).Table("issue_labels").
				Select("issue_id").
				Where("label_id IN ?", filter.LabelIDs).
				Group("issue_id").
				Having("COUNT(DISTINCT label_id) = ?", len(filter.LabelIDs))
			query = query.Where("id IN (?)", sub)
		} else {
			sub := query.Session(&gorm.Session{}).Table("issue_labels").
				Select("issue_id").
				Where("label_id IN ?", filter.LabelIDs)
			query = query.Where("id IN (?)", sub)
		}
	}
	return query
}

func strPtr(v string) *string {
	value := v
	return &value
}

func strPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func isUniqueViolationConstraint(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code != "23505" {
			return false
		}
		if constraint == "" {
			return true
		}
		return pgErr.ConstraintName == constraint
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if string(pqErr.Code) != "23505" {
			return false
		}
		if constraint == "" {
			return true
		}
		return pqErr.Constraint == constraint
	}

	return false
}
