package services

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	cachepkg "github.com/abhinavmaity/linear-lite/backend/internal/cache"
	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
	"github.com/google/uuid"
)

const (
	maxSprintNameLength        = 255
	maxSprintDescriptionLength = 10000
)

type SprintListInput struct {
	Page      int
	Limit     int
	ProjectID *string
	Status    *string
	Search    string
	SortBy    string
	SortOrder string
}

type SprintCreateInput struct {
	Name        string
	Description *string
	ProjectID   string
	StartDate   string
	EndDate     string
	Status      *string
}

type SprintUpdateInput struct {
	Name        *string
	Description **string
	StartDate   *string
	EndDate     *string
	Status      *string
}

type SprintService struct {
	repo     repositories.SprintRepository
	projects repositories.ProjectRepository
	cache    *cachepkg.Store
}

func NewSprintService(
	repo repositories.SprintRepository,
	projects repositories.ProjectRepository,
	cache *cachepkg.Store,
) *SprintService {
	return &SprintService{
		repo:     repo,
		projects: projects,
		cache:    cache,
	}
}

func (s *SprintService) List(ctx context.Context, input SprintListInput) ([]SprintSummary, int64, *apperrors.AppError) {
	projectID := ""
	if input.ProjectID != nil {
		projectID = *input.ProjectID
	}
	status := ""
	if input.Status != nil {
		status = *input.Status
	}
	cacheKey := buildListCacheKey(
		"sprints",
		intToCachePart(input.Page),
		intToCachePart(input.Limit),
		projectID,
		status,
		input.Search,
		input.SortBy,
		input.SortOrder,
	)
	if s.cache != nil {
		var cached struct {
			Items []SprintSummary `json:"items"`
			Total int64           `json:"total"`
		}
		if found, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && found {
			return cached.Items, cached.Total, nil
		}
	}

	sprints, total, err := s.repo.List(ctx, repositories.SprintListFilter{
		PaginationInput: repositories.PaginationInput{
			Page:  input.Page,
			Limit: input.Limit,
		},
		SortInput: repositories.SortInput{
			By:    input.SortBy,
			Order: input.SortOrder,
		},
		Search:    input.Search,
		ProjectID: input.ProjectID,
		Status:    input.Status,
	})
	if err != nil {
		return nil, 0, apperrors.Internal("failed to list sprints")
	}

	items := make([]SprintSummary, 0, len(sprints))
	for _, sprint := range sprints {
		items = append(items, SprintSummary{
			ID:          sprint.ID,
			Name:        sprint.Name,
			Description: sprint.Description,
			ProjectID:   sprint.ProjectID,
			StartDate:   sprint.StartDate.Format(time.DateOnly),
			EndDate:     sprint.EndDate.Format(time.DateOnly),
			Status:      sprint.Status,
			CreatedAt:   sprint.CreatedAt,
			UpdatedAt:   sprint.UpdatedAt,
			IssueCounts: mapIssueCounts(sprint.IssueCounts),
		})
	}

	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, struct {
			Items []SprintSummary `json:"items"`
			Total int64           `json:"total"`
		}{
			Items: items,
			Total: total,
		}, 2*time.Minute)
	}

	return items, total, nil
}

func (s *SprintService) Create(ctx context.Context, input SprintCreateInput) (*SprintDetail, *apperrors.AppError) {
	fields := apperrors.FieldErrors{}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		fields["name"] = "is required"
	} else if utf8.RuneCountInString(name) > maxSprintNameLength {
		fields["name"] = "must be less than or equal to 255 characters"
	}

	description := normalizeSprintOptional(input.Description)
	if description != nil && utf8.RuneCountInString(*description) > maxSprintDescriptionLength {
		fields["description"] = "must be less than or equal to 10000 characters"
	}

	projectID := strings.TrimSpace(input.ProjectID)
	if _, err := uuid.Parse(projectID); err != nil {
		fields["project_id"] = "must be a valid UUID"
	}

	startDate, appErr := parseDateField("start_date", input.StartDate)
	if appErr != nil {
		fields["start_date"] = appErr.Fields["start_date"]
	}
	endDate, appErr := parseDateField("end_date", input.EndDate)
	if appErr != nil {
		fields["end_date"] = appErr.Fields["end_date"]
	}

	status := models.SprintStatusPlanned
	if input.Status != nil {
		status = strings.TrimSpace(*input.Status)
		if enumErr := validateSprintStatus(status); enumErr != nil {
			fields["status"] = enumErr.Fields["status"]
		}
	}

	if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(startDate) {
		fields["end_date"] = "must be greater than or equal to start_date"
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("one or more fields are invalid", fields)
	}

	projectExists, err := s.projects.ExistsByID(ctx, projectID)
	if err != nil {
		return nil, apperrors.Internal("failed to validate project")
	}
	if !projectExists {
		return nil, apperrors.NotFound("project not found")
	}

	sprint := &models.Sprint{
		Name:        name,
		Description: description,
		ProjectID:   projectID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      status,
	}
	if err := s.repo.Create(ctx, sprint); err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Conflict("only one active sprint is allowed per project", apperrors.FieldErrors{
				"status": "an active sprint already exists for this project",
			})
		}
		return nil, apperrors.Internal("failed to create sprint")
	}

	s.invalidateSprintCaches(ctx)

	return s.Get(ctx, sprint.ID)
}

func (s *SprintService) Get(ctx context.Context, id string) (*SprintDetail, *apperrors.AppError) {
	cacheKey := buildDetailCacheKey("sprints", id)
	if s.cache != nil {
		var cached SprintDetail
		if found, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && found {
			return &cached, nil
		}
	}

	sprintMap, err := s.repo.SummariesByIDs(ctx, []string{id})
	if err != nil {
		return nil, apperrors.Internal("failed to load sprint")
	}
	sprintRow, ok := sprintMap[id]
	if !ok {
		return nil, apperrors.NotFound("sprint not found")
	}

	projectMap, err := s.projects.SummariesByIDs(ctx, []string{sprintRow.ProjectID})
	if err != nil {
		return nil, apperrors.Internal("failed to load sprint project")
	}
	projectRow, ok := projectMap[sprintRow.ProjectID]
	if !ok {
		return nil, apperrors.NotFound("project not found")
	}

	detail := SprintDetail{
		SprintSummary: *mapSprintRow(sprintRow),
		Project: ProjectSummary{
			ID:          projectRow.ID,
			Name:        projectRow.Name,
			Description: projectRow.Description,
			Key:         projectRow.Key,
			CreatedBy:   projectRow.CreatedBy,
			CreatedAt:   projectRow.CreatedAt,
			UpdatedAt:   projectRow.UpdatedAt,
			IssueCounts: mapIssueCounts(projectRow.IssueCounts),
		},
	}
	if projectRow.ActiveSprint != nil {
		detail.Project.ActiveSprint = mapSprintRow(*projectRow.ActiveSprint)
	}

	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, detail, 2*time.Minute)
	}
	return &detail, nil
}

func (s *SprintService) Update(ctx context.Context, id string, input SprintUpdateInput) (*SprintDetail, *apperrors.AppError) {
	sprint, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("sprint not found")
		}
		return nil, apperrors.Internal("failed to load sprint")
	}

	fields := apperrors.FieldErrors{}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			fields["name"] = "is required"
		} else if utf8.RuneCountInString(name) > maxSprintNameLength {
			fields["name"] = "must be less than or equal to 255 characters"
		} else {
			sprint.Name = name
		}
	}
	if input.Description != nil {
		description := normalizeSprintOptional(*input.Description)
		if description != nil && utf8.RuneCountInString(*description) > maxSprintDescriptionLength {
			fields["description"] = "must be less than or equal to 10000 characters"
		} else {
			sprint.Description = description
		}
	}
	if input.Status != nil {
		status := strings.TrimSpace(*input.Status)
		if enumErr := validateSprintStatus(status); enumErr != nil {
			fields["status"] = enumErr.Fields["status"]
		} else {
			sprint.Status = status
		}
	}

	startDate := sprint.StartDate
	endDate := sprint.EndDate
	if input.StartDate != nil {
		parsed, appErr := parseDateField("start_date", *input.StartDate)
		if appErr != nil {
			fields["start_date"] = appErr.Fields["start_date"]
		} else {
			startDate = parsed
		}
	}
	if input.EndDate != nil {
		parsed, appErr := parseDateField("end_date", *input.EndDate)
		if appErr != nil {
			fields["end_date"] = appErr.Fields["end_date"]
		} else {
			endDate = parsed
		}
	}
	if endDate.Before(startDate) {
		fields["end_date"] = "must be greater than or equal to start_date"
	}
	sprint.StartDate = startDate
	sprint.EndDate = endDate

	if len(fields) > 0 {
		return nil, apperrors.Validation("one or more fields are invalid", fields)
	}

	if err := s.repo.Update(ctx, sprint); err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Conflict("only one active sprint is allowed per project", apperrors.FieldErrors{
				"status": "an active sprint already exists for this project",
			})
		}
		return nil, apperrors.Internal("failed to update sprint")
	}

	s.invalidateSprintCaches(ctx)

	return s.Get(ctx, id)
}

func (s *SprintService) Delete(ctx context.Context, id string) *apperrors.AppError {
	sprint, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return apperrors.NotFound("sprint not found")
		}
		return apperrors.Internal("failed to load sprint")
	}

	if sprint.Status == models.SprintStatusActive {
		return apperrors.Conflict("active sprint cannot be deleted", apperrors.FieldErrors{
			"id": "sprint must not be active",
		})
	}

	issueCount, err := s.repo.CountIssuesBySprintID(ctx, id)
	if err != nil {
		return apperrors.Internal("failed to validate sprint deletion")
	}
	if issueCount > 0 {
		return apperrors.Conflict("sprint cannot be deleted while issues reference it", apperrors.FieldErrors{
			"id": "sprint has dependent issues",
		})
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return apperrors.NotFound("sprint not found")
		}
		return apperrors.Internal("failed to delete sprint")
	}

	s.invalidateSprintCaches(ctx)

	return nil
}

func (s *SprintService) invalidateSprintCaches(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.DeleteByPrefix(ctx, "sprints:")
	_ = s.cache.DeleteByPrefix(ctx, "projects:")
	_ = s.cache.DeleteByPrefix(ctx, "dashboard:")
}

func validateSprintStatus(status string) *apperrors.AppError {
	allowed := []string{models.SprintStatusPlanned, models.SprintStatusActive, models.SprintStatusCompleted}
	for _, option := range allowed {
		if status == option {
			return nil
		}
	}
	return apperrors.Validation("one or more fields are invalid", apperrors.FieldErrors{
		"status": "must be one of: planned, active, completed",
	})
}

func parseDateField(field, raw string) (time.Time, *apperrors.AppError) {
	parsed, err := time.Parse(time.DateOnly, strings.TrimSpace(raw))
	if err != nil {
		return time.Time{}, apperrors.Validation("one or more fields are invalid", apperrors.FieldErrors{
			field: "must use YYYY-MM-DD format",
		})
	}
	return parsed, nil
}

func normalizeSprintOptional(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}
