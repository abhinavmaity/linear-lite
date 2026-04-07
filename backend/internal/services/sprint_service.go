package services

import (
	"context"
	"time"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
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

type SprintService struct {
	repo repositories.SprintRepository
}

func NewSprintService(repo repositories.SprintRepository) *SprintService {
	return &SprintService{repo: repo}
}

func (s *SprintService) List(ctx context.Context, input SprintListInput) ([]SprintSummary, int64, *apperrors.AppError) {
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

	return items, total, nil
}
