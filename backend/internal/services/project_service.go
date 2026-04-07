package services

import (
	"context"
	"time"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
)

type ProjectListInput struct {
	Page      int
	Limit     int
	Search    string
	SortBy    string
	SortOrder string
}

type ProjectService struct {
	repo repositories.ProjectRepository
}

func NewProjectService(repo repositories.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
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
