package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SprintRepositoryDB struct {
	db *gorm.DB
}

func NewSprintRepository(db *gorm.DB) *SprintRepositoryDB {
	return &SprintRepositoryDB{db: db}
}

func (r *SprintRepositoryDB) FindByID(ctx context.Context, id string) (*models.Sprint, error) {
	var sprint models.Sprint
	query := r.db.WithContext(ctx).Where("id = ?", strings.TrimSpace(id)).First(&sprint)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if query.Error != nil {
		return nil, query.Error
	}
	return &sprint, nil
}

func (r *SprintRepositoryDB) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Sprint{}).
		Where("id = ?", strings.TrimSpace(id)).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SprintRepositoryDB) List(ctx context.Context, filter SprintListFilter) ([]SprintSummaryRow, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Sprint{})

	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where("name ILIKE ?", pattern)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []SprintSummaryRow{}, 0, nil
	}

	sortBy := filter.By
	if sortBy == "" {
		sortBy = "start_date"
	}
	order := strings.ToLower(strings.TrimSpace(filter.Order))
	if order != "asc" {
		order = "desc"
	}

	orderClause := "start_date desc"
	switch sortBy {
	case "name":
		orderClause = "name " + order
	case "start_date":
		orderClause = "start_date " + order
	case "end_date":
		orderClause = "end_date " + order
	case "created_at":
		orderClause = "created_at " + order
	}

	var sprints []models.Sprint
	err := query.Order(orderClause).
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Find(&sprints).Error
	if err != nil {
		return nil, 0, err
	}

	sprintIDs := make([]string, 0, len(sprints))
	for _, sprint := range sprints {
		sprintIDs = append(sprintIDs, sprint.ID)
	}

	countsBySprint, err := r.loadIssueCountsBySprint(ctx, sprintIDs)
	if err != nil {
		return nil, 0, err
	}

	rows := make([]SprintSummaryRow, 0, len(sprints))
	for _, sprint := range sprints {
		rows = append(rows, SprintSummaryRow{
			ID:          sprint.ID,
			Name:        sprint.Name,
			Description: sprint.Description,
			ProjectID:   sprint.ProjectID,
			StartDate:   sprint.StartDate,
			EndDate:     sprint.EndDate,
			Status:      sprint.Status,
			CreatedAt:   sprint.CreatedAt,
			UpdatedAt:   sprint.UpdatedAt,
			IssueCounts: countsBySprint[sprint.ID],
		})
	}

	return rows, total, nil
}

func (r *SprintRepositoryDB) SummariesByIDs(ctx context.Context, ids []string) (map[string]SprintSummaryRow, error) {
	out := make(map[string]SprintSummaryRow, len(ids))
	if len(ids) == 0 {
		return out, nil
	}

	var sprints []models.Sprint
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&sprints).Error
	if err != nil {
		return nil, err
	}

	sprintIDs := make([]string, 0, len(sprints))
	for _, sprint := range sprints {
		sprintIDs = append(sprintIDs, sprint.ID)
	}
	countsBySprint, err := r.loadIssueCountsBySprint(ctx, sprintIDs)
	if err != nil {
		return nil, err
	}

	for _, sprint := range sprints {
		out[sprint.ID] = SprintSummaryRow{
			ID:          sprint.ID,
			Name:        sprint.Name,
			Description: sprint.Description,
			ProjectID:   sprint.ProjectID,
			StartDate:   sprint.StartDate,
			EndDate:     sprint.EndDate,
			Status:      sprint.Status,
			CreatedAt:   sprint.CreatedAt,
			UpdatedAt:   sprint.UpdatedAt,
			IssueCounts: countsBySprint[sprint.ID],
		}
	}

	return out, nil
}

func (r *SprintRepositoryDB) Create(ctx context.Context, sprint *models.Sprint) error {
	if err := r.db.WithContext(ctx).Create(sprint).Error; err != nil {
		if isSprintUniqueViolation(err, "uq_sprints_one_active_per_project") {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *SprintRepositoryDB) Update(ctx context.Context, sprint *models.Sprint) error {
	if err := r.db.WithContext(ctx).Save(sprint).Error; err != nil {
		if isSprintUniqueViolation(err, "uq_sprints_one_active_per_project") {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *SprintRepositoryDB) Delete(ctx context.Context, id string) error {
	query := r.db.WithContext(ctx).Delete(&models.Sprint{}, "id = ?", strings.TrimSpace(id))
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *SprintRepositoryDB) CountIssuesBySprintID(ctx context.Context, id string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Issue{}).
		Where("sprint_id = ?", strings.TrimSpace(id)).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *SprintRepositoryDB) loadIssueCountsBySprint(ctx context.Context, sprintIDs []string) (map[string]IssueCounts, error) {
	out := make(map[string]IssueCounts, len(sprintIDs))
	if len(sprintIDs) == 0 {
		return out, nil
	}

	var rows []issueCountsRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			sprint_id AS resource_id,
			COUNT(*)::int AS total,
			COUNT(*) FILTER (WHERE status = 'backlog')::int AS backlog,
			COUNT(*) FILTER (WHERE status = 'todo')::int AS todo,
			COUNT(*) FILTER (WHERE status = 'in_progress')::int AS in_progress,
			COUNT(*) FILTER (WHERE status = 'in_review')::int AS in_review,
			COUNT(*) FILTER (WHERE status = 'done')::int AS done,
			COUNT(*) FILTER (WHERE status = 'cancelled')::int AS cancelled
		FROM issues
		WHERE archived_at IS NULL
		  AND sprint_id IN ?
		GROUP BY sprint_id
	`, sprintIDs).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		out[row.ResourceID] = IssueCounts{
			Total:      row.Total,
			Backlog:    row.Backlog,
			Todo:       row.Todo,
			InProgress: row.InProgress,
			InReview:   row.InReview,
			Done:       row.Done,
			Cancelled:  row.Cancelled,
		}
	}
	return out, nil
}

func isSprintUniqueViolation(err error, constraint string) bool {
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
