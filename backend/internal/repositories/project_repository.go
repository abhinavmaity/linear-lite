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

type ProjectRepositoryDB struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepositoryDB {
	return &ProjectRepositoryDB{db: db}
}

func (r *ProjectRepositoryDB) FindByID(ctx context.Context, id string) (*models.Project, error) {
	var project models.Project
	query := r.db.WithContext(ctx).Where("id = ?", strings.TrimSpace(id)).First(&project)
	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, query.Error
	}
	return &project, nil
}

func (r *ProjectRepositoryDB) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Project{}).
		Where("id = ?", strings.TrimSpace(id)).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProjectRepositoryDB) List(ctx context.Context, filter ProjectListFilter) ([]ProjectSummaryRow, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Project{})

	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR key ILIKE ?", pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []ProjectSummaryRow{}, 0, nil
	}

	sortBy := filter.By
	if sortBy == "" {
		sortBy = "name"
	}
	order := strings.ToLower(strings.TrimSpace(filter.Order))
	if order != "desc" {
		order = "asc"
	}

	orderClause := "name asc"
	switch sortBy {
	case "name":
		orderClause = "name " + order
	case "created_at":
		orderClause = "created_at " + order
	case "updated_at":
		orderClause = "updated_at " + order
	}

	var projects []models.Project
	err := query.Order(orderClause).
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	projectIDs := make([]string, 0, len(projects))
	for _, p := range projects {
		projectIDs = append(projectIDs, p.ID)
	}

	projectCounts, err := r.loadIssueCountsByProject(ctx, projectIDs)
	if err != nil {
		return nil, 0, err
	}
	activeSprints, err := r.loadActiveSprintsByProject(ctx, projectIDs)
	if err != nil {
		return nil, 0, err
	}
	sprintCounts, err := r.loadIssueCountsBySprint(ctx, activeSprints)
	if err != nil {
		return nil, 0, err
	}

	rows := make([]ProjectSummaryRow, 0, len(projects))
	for _, p := range projects {
		row := ProjectSummaryRow{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Key:         p.Key,
			CreatedBy:   p.CreatedBy,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			IssueCounts: projectCounts[p.ID],
		}

		if active, ok := activeSprints[p.ID]; ok {
			active.IssueCounts = sprintCounts[active.ID]
			row.ActiveSprint = active
		}

		rows = append(rows, row)
	}

	return rows, total, nil
}

func (r *ProjectRepositoryDB) SummariesByIDs(ctx context.Context, ids []string) (map[string]ProjectSummaryRow, error) {
	out := make(map[string]ProjectSummaryRow, len(ids))
	if len(ids) == 0 {
		return out, nil
	}

	var projects []models.Project
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&projects).Error; err != nil {
		return nil, err
	}

	projectIDs := make([]string, 0, len(projects))
	for _, p := range projects {
		projectIDs = append(projectIDs, p.ID)
	}
	projectCounts, err := r.loadIssueCountsByProject(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	activeSprints, err := r.loadActiveSprintsByProject(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	sprintCounts, err := r.loadIssueCountsBySprint(ctx, activeSprints)
	if err != nil {
		return nil, err
	}

	for _, p := range projects {
		row := ProjectSummaryRow{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			Key:          p.Key,
			CreatedBy:    p.CreatedBy,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
			IssueCounts:  projectCounts[p.ID],
			ActiveSprint: nil,
		}
		if active, ok := activeSprints[p.ID]; ok {
			active.IssueCounts = sprintCounts[active.ID]
			row.ActiveSprint = active
		}
		out[p.ID] = row
	}

	return out, nil
}

func (r *ProjectRepositoryDB) ListSprintsByProjectID(ctx context.Context, projectID string) ([]SprintSummaryRow, error) {
	var sprints []models.Sprint
	err := r.db.WithContext(ctx).
		Where("project_id = ?", strings.TrimSpace(projectID)).
		Order("start_date DESC").
		Find(&sprints).Error
	if err != nil {
		return nil, err
	}
	if len(sprints) == 0 {
		return []SprintSummaryRow{}, nil
	}

	sprintIDs := make([]string, 0, len(sprints))
	for _, sprint := range sprints {
		sprintIDs = append(sprintIDs, sprint.ID)
	}
	countsBySprint, err := r.loadIssueCountsBySprintIDs(ctx, sprintIDs)
	if err != nil {
		return nil, err
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
	return rows, nil
}

func (r *ProjectRepositoryDB) Create(ctx context.Context, project *models.Project) error {
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		if isProjectUniqueViolation(err, "uq_projects_key") {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *ProjectRepositoryDB) Update(ctx context.Context, project *models.Project) error {
	if err := r.db.WithContext(ctx).Save(project).Error; err != nil {
		if isProjectUniqueViolation(err, "uq_projects_key") {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *ProjectRepositoryDB) Delete(ctx context.Context, id string) error {
	query := r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", strings.TrimSpace(id))
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProjectRepositoryDB) CountIssuesByProjectID(ctx context.Context, id string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Issue{}).
		Where("project_id = ?", strings.TrimSpace(id)).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProjectRepositoryDB) CountSprintsByProjectID(ctx context.Context, id string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Sprint{}).
		Where("project_id = ?", strings.TrimSpace(id)).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

type issueCountsRow struct {
	ResourceID string `gorm:"column:resource_id"`
	Total      int    `gorm:"column:total"`
	Backlog    int    `gorm:"column:backlog"`
	Todo       int    `gorm:"column:todo"`
	InProgress int    `gorm:"column:in_progress"`
	InReview   int    `gorm:"column:in_review"`
	Done       int    `gorm:"column:done"`
	Cancelled  int    `gorm:"column:cancelled"`
}

func (r *ProjectRepositoryDB) loadIssueCountsByProject(ctx context.Context, projectIDs []string) (map[string]IssueCounts, error) {
	out := make(map[string]IssueCounts, len(projectIDs))
	if len(projectIDs) == 0 {
		return out, nil
	}

	var rows []issueCountsRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			project_id AS resource_id,
			COUNT(*)::int AS total,
			COUNT(*) FILTER (WHERE status = 'backlog')::int AS backlog,
			COUNT(*) FILTER (WHERE status = 'todo')::int AS todo,
			COUNT(*) FILTER (WHERE status = 'in_progress')::int AS in_progress,
			COUNT(*) FILTER (WHERE status = 'in_review')::int AS in_review,
			COUNT(*) FILTER (WHERE status = 'done')::int AS done,
			COUNT(*) FILTER (WHERE status = 'cancelled')::int AS cancelled
		FROM issues
		WHERE archived_at IS NULL
		  AND project_id IN ?
		GROUP BY project_id
	`, projectIDs).Scan(&rows).Error
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

func (r *ProjectRepositoryDB) loadActiveSprintsByProject(ctx context.Context, projectIDs []string) (map[string]*SprintSummaryRow, error) {
	out := make(map[string]*SprintSummaryRow, len(projectIDs))
	if len(projectIDs) == 0 {
		return out, nil
	}

	var sprints []models.Sprint
	err := r.db.WithContext(ctx).Model(&models.Sprint{}).
		Where("status = ?", models.SprintStatusActive).
		Where("project_id IN ?", projectIDs).
		Find(&sprints).Error
	if err != nil {
		return nil, err
	}

	for _, s := range sprints {
		sprint := &SprintSummaryRow{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			ProjectID:   s.ProjectID,
			StartDate:   s.StartDate,
			EndDate:     s.EndDate,
			Status:      s.Status,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		}
		out[s.ProjectID] = sprint
	}

	return out, nil
}

func (r *ProjectRepositoryDB) loadIssueCountsBySprint(ctx context.Context, activeSprints map[string]*SprintSummaryRow) (map[string]IssueCounts, error) {
	sprintIDs := make([]string, 0, len(activeSprints))
	for _, sprint := range activeSprints {
		sprintIDs = append(sprintIDs, sprint.ID)
	}
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

func (r *ProjectRepositoryDB) loadIssueCountsBySprintIDs(ctx context.Context, sprintIDs []string) (map[string]IssueCounts, error) {
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

func isProjectUniqueViolation(err error, constraint string) bool {
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
