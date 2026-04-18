package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	cachepkg "github.com/abhinavmaity/linear-lite/backend/internal/cache"
	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
)

const (
	maxProjectNameLength        = 255
	maxProjectDescriptionLength = 10000
)

var projectKeyPattern = regexp.MustCompile(`^[A-Z0-9]{2,10}$`)

type ProjectListInput struct {
	Page      int
	Limit     int
	Search    string
	SortBy    string
	SortOrder string
}

type ProjectCreateInput struct {
	Name        string
	Description *string
	Key         string
}

type ProjectUpdateInput struct {
	Name        *string
	Description **string
	Key         *string
}

type ProjectService struct {
	repo  repositories.ProjectRepository
	users repositories.UserReadRepository
	cache *cachepkg.Store
}

func NewProjectService(
	repo repositories.ProjectRepository,
	users repositories.UserReadRepository,
	cache *cachepkg.Store,
) *ProjectService {
	return &ProjectService{
		repo:  repo,
		users: users,
		cache: cache,
	}
}

func (s *ProjectService) List(ctx context.Context, input ProjectListInput) ([]ProjectSummary, int64, *apperrors.AppError) {
	cacheKey := buildListCacheKey(
		"projects",
		intToCachePart(input.Page),
		intToCachePart(input.Limit),
		input.Search,
		input.SortBy,
		input.SortOrder,
	)
	if s.cache != nil {
		var cached struct {
			Items []ProjectSummary `json:"items"`
			Total int64            `json:"total"`
		}
		if found, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && found {
			return cached.Items, cached.Total, nil
		}
	}

	projects, total, err := s.repo.List(ctx, repositories.ProjectListFilter{
		PaginationInput: repositories.PaginationInput{
			Page:  input.Page,
			Limit: input.Limit,
		},
		SortInput: repositories.SortInput{
			By:    input.SortBy,
			Order: input.SortOrder,
		},
		Search: input.Search,
	})
	if err != nil {
		return nil, 0, apperrors.Internal("failed to list projects")
	}

	items := make([]ProjectSummary, 0, len(projects))
	for _, project := range projects {
		item := ProjectSummary{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			Key:         project.Key,
			CreatedBy:   project.CreatedBy,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
			IssueCounts: mapIssueCounts(project.IssueCounts),
		}
		if project.ActiveSprint != nil {
			item.ActiveSprint = mapSprintRow(*project.ActiveSprint)
		}
		items = append(items, item)
	}

	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, struct {
			Items []ProjectSummary `json:"items"`
			Total int64            `json:"total"`
		}{
			Items: items,
			Total: total,
		}, 2*time.Minute)
	}

	return items, total, nil
}

func mapIssueCounts(counts repositories.IssueCounts) IssueCounts {
	return IssueCounts{
		Total:      counts.Total,
		Backlog:    counts.Backlog,
		Todo:       counts.Todo,
		InProgress: counts.InProgress,
		InReview:   counts.InReview,
		Done:       counts.Done,
		Cancelled:  counts.Cancelled,
	}
}

func mapSprintRow(row repositories.SprintSummaryRow) *SprintSummary {
	return &SprintSummary{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		ProjectID:   row.ProjectID,
		StartDate:   row.StartDate.Format(time.DateOnly),
		EndDate:     row.EndDate.Format(time.DateOnly),
		Status:      row.Status,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		IssueCounts: mapIssueCounts(row.IssueCounts),
	}
}

func (s *ProjectService) Create(ctx context.Context, actorID string, input ProjectCreateInput) (*ProjectDetail, *apperrors.AppError) {
	fields := apperrors.FieldErrors{}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		fields["name"] = "Project name is required."
	} else if utf8.RuneCountInString(name) > maxProjectNameLength {
		fields["name"] = "Project name must be 255 characters or fewer."
	}

	description := normalizeProjectOptional(input.Description)
	if description != nil && utf8.RuneCountInString(*description) > maxProjectDescriptionLength {
		fields["description"] = "Description must be 10000 characters or fewer."
	}

	key := strings.TrimSpace(input.Key)
	if key == "" {
		fields["key"] = "Project key is required."
	} else if !projectKeyPattern.MatchString(key) {
		fields["key"] = "Project key must be 2-10 uppercase letters or numbers."
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", fields)
	}

	project := &models.Project{
		Name:        name,
		Description: description,
		Key:         key,
		CreatedBy:   actorID,
	}
	if err := s.repo.Create(ctx, project); err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Conflict("project key already exists", apperrors.FieldErrors{
				"key": "This project key is already in use.",
			})
		}
		return nil, apperrors.Internal("failed to create project")
	}

	s.invalidateProjectCaches(ctx)

	return s.Get(ctx, project.ID)
}

func (s *ProjectService) Get(ctx context.Context, id string) (*ProjectDetail, *apperrors.AppError) {
	cacheKey := buildDetailCacheKey("projects", id)
	if s.cache != nil {
		var cached ProjectDetail
		if found, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && found {
			return &cached, nil
		}
	}

	summaryMap, err := s.repo.SummariesByIDs(ctx, []string{id})
	if err != nil {
		return nil, apperrors.Internal("failed to load project")
	}
	projectRow, ok := summaryMap[id]
	if !ok {
		return nil, apperrors.NotFound("project not found")
	}

	creator, err := s.users.FindByID(ctx, projectRow.CreatedBy)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("project creator not found")
		}
		return nil, apperrors.Internal("failed to load project creator")
	}

	sprintRows, err := s.repo.ListSprintsByProjectID(ctx, id)
	if err != nil {
		return nil, apperrors.Internal("failed to load project sprints")
	}

	sprints := make([]SprintSummary, 0, len(sprintRows))
	for _, row := range sprintRows {
		sprints = append(sprints, *mapSprintRow(row))
	}

	detail := ProjectDetail{
		ProjectSummary: ProjectSummary{
			ID:          projectRow.ID,
			Name:        projectRow.Name,
			Description: projectRow.Description,
			Key:         projectRow.Key,
			CreatedBy:   projectRow.CreatedBy,
			CreatedAt:   projectRow.CreatedAt,
			UpdatedAt:   projectRow.UpdatedAt,
			IssueCounts: mapIssueCounts(projectRow.IssueCounts),
		},
		Creator: UserSummary{
			ID:        creator.ID,
			Email:     creator.Email,
			Name:      creator.Name,
			AvatarURL: creator.AvatarURL,
			CreatedAt: creator.CreatedAt,
			UpdatedAt: creator.UpdatedAt,
		},
		Sprints: sprints,
	}

	if projectRow.ActiveSprint != nil {
		detail.ActiveSprint = mapSprintRow(*projectRow.ActiveSprint)
	}

	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, detail, 2*time.Minute)
	}

	return &detail, nil
}

func (s *ProjectService) Update(ctx context.Context, id string, input ProjectUpdateInput) (*ProjectDetail, *apperrors.AppError) {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("project not found")
		}
		return nil, apperrors.Internal("failed to load project")
	}

	fields := apperrors.FieldErrors{}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			fields["name"] = "Project name is required."
		} else if utf8.RuneCountInString(name) > maxProjectNameLength {
			fields["name"] = "Project name must be 255 characters or fewer."
		} else {
			project.Name = name
		}
	}
	if input.Description != nil {
		description := normalizeProjectOptional(*input.Description)
		if description != nil && utf8.RuneCountInString(*description) > maxProjectDescriptionLength {
			fields["description"] = "Description must be 10000 characters or fewer."
		} else {
			project.Description = description
		}
	}
	if input.Key != nil {
		key := strings.TrimSpace(*input.Key)
		if key == "" {
			fields["key"] = "Project key is required."
		} else if !projectKeyPattern.MatchString(key) {
			fields["key"] = "Project key must be 2-10 uppercase letters or numbers."
		} else if key != project.Key {
			issueCount, countErr := s.repo.CountIssuesByProjectID(ctx, id)
			if countErr != nil {
				return nil, apperrors.Internal("failed to validate project key update")
			}
			if issueCount > 0 {
				return nil, apperrors.Conflict("project key cannot change once issues exist", apperrors.FieldErrors{
					"key": "Project key cannot be changed while the project has issues.",
				})
			}
			project.Key = key
		}
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", fields)
	}

	if err := s.repo.Update(ctx, project); err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Conflict("project key already exists", apperrors.FieldErrors{
				"key": "This project key is already in use.",
			})
		}
		return nil, apperrors.Internal("failed to update project")
	}

	s.invalidateProjectCaches(ctx)

	return s.Get(ctx, id)
}

func (s *ProjectService) Delete(ctx context.Context, id string) *apperrors.AppError {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return apperrors.NotFound("project not found")
		}
		return apperrors.Internal("failed to load project")
	}

	issueCount, err := s.repo.CountIssuesByProjectID(ctx, id)
	if err != nil {
		return apperrors.Internal("failed to validate project deletion")
	}
	if issueCount > 0 {
		return apperrors.Conflict("project cannot be deleted while issues exist", apperrors.FieldErrors{
			"id": "project has dependent issues",
		})
	}

	sprintCount, err := s.repo.CountSprintsByProjectID(ctx, id)
	if err != nil {
		return apperrors.Internal("failed to validate project deletion")
	}
	if sprintCount > 0 {
		return apperrors.Conflict("project cannot be deleted while sprints exist", apperrors.FieldErrors{
			"id": "project has dependent sprints",
		})
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return apperrors.NotFound("project not found")
		}
		return apperrors.Internal("failed to delete project")
	}

	s.invalidateProjectCaches(ctx)

	return nil
}

func (s *ProjectService) invalidateProjectCaches(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.DeleteByPrefix(ctx, "projects:")
	_ = s.cache.DeleteByPrefix(ctx, "dashboard:")
}

func normalizeProjectOptional(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}
