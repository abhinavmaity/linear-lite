package services

import (
	"context"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
)

type LabelListInput struct {
	Page      int
	Limit     int
	Search    string
	SortBy    string
	SortOrder string
}

type LabelService struct {
	repo repositories.LabelRepository
}

func NewLabelService(repo repositories.LabelRepository) *LabelService {
	return &LabelService{repo: repo}
}

func (s *LabelService) List(ctx context.Context, input LabelListInput) ([]LabelSummary, int64, *apperrors.AppError) {
	labels, total, err := s.repo.List(ctx, repositories.LabelListFilter{
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
		return nil, 0, apperrors.Internal("failed to list labels")
	}

	items := make([]LabelSummary, 0, len(labels))
	for _, label := range labels {
		items = append(items, LabelSummary{
			ID:          label.ID,
			Name:        label.Name,
			Color:       label.Color,
			Description: label.Description,
			CreatedAt:   label.CreatedAt,
			UpdatedAt:   label.UpdatedAt,
		})
	}

	return items, total, nil
}
