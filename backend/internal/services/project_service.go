package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

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
}

func NewProjectService(
	repo repositories.ProjectRepository,
	users repositories.UserReadRepository,
) *ProjectService {
	return &ProjectService{
		repo:  repo,
		users: users,
	}
}

func (s *ProjectService) List(ctx context.Context, input ProjectListInput) ([]ProjectSummary, int64, *apperrors.AppError) {
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
		fields["name"] = "is required"
	} else if utf8.RuneCountInString(name) > maxProjectNameLength {
		fields["name"] = "must be less than or equal to 255 characters"
	}

	description := normalizeProjectOptional(input.Description)
	if description != nil && utf8.RuneCountInString(*description) > maxProjectDescriptionLength {
		fields["description"] = "must be less than or equal to 10000 characters"
	}

	key := strings.TrimSpace(input.Key)
	if key == "" {
		fields["key"] = "is required"
	} else if !projectKeyPattern.MatchString(key) {
		fields["key"] = "must match ^[A-Z0-9]{2,10}$"
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("one or more fields are invalid", fields)
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
				"key": "already in use",
			})
		}
		return nil, apperrors.Internal("failed to create project")
	}

	return s.Get(ctx, project.ID)
}

func (s *ProjectService) Get(ctx context.Context, id string) (*ProjectDetail, *apperrors.AppError) {
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
			fields["name"] = "is required"
		} else if utf8.RuneCountInString(name) > maxProjectNameLength {
			fields["name"] = "must be less than or equal to 255 characters"
		} else {
			project.Name = name
		}
	}
	if input.Description != nil {
		description := normalizeProjectOptional(*input.Description)
		if description != nil && utf8.RuneCountInString(*description) > maxProjectDescriptionLength {
			fields["description"] = "must be less than or equal to 10000 characters"
		} else {
			project.Description = description
		}
	}
	if input.Key != nil {
		key := strings.TrimSpace(*input.Key)
		if key == "" {
			fields["key"] = "is required"
		} else if !projectKeyPattern.MatchString(key) {
			fields["key"] = "must match ^[A-Z0-9]{2,10}$"
		} else if key != project.Key {
			issueCount, countErr := s.repo.CountIssuesByProjectID(ctx, id)
			if countErr != nil {
				return nil, apperrors.Internal("failed to validate project key update")
			}
			if issueCount > 0 {
				return nil, apperrors.Conflict("project key cannot change once issues exist", apperrors.FieldErrors{
					"key": "cannot be changed when project has issues",
				})
			}
			project.Key = key
		}
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("one or more fields are invalid", fields)
	}

	if err := s.repo.Update(ctx, project); err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Conflict("project key already exists", apperrors.FieldErrors{
				"key": "already in use",
			})
		}
		return nil, apperrors.Internal("failed to update project")
	}

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

	return nil
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
