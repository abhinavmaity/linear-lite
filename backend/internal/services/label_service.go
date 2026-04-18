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
	maxLabelNameLength        = 50
	maxLabelDescriptionLength = 1000
)

var labelColorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type LabelListInput struct {
	Page      int
	Limit     int
	Search    string
	SortBy    string
	SortOrder string
}

type LabelCreateInput struct {
	Name        string
	Color       string
	Description *string
}

type LabelUpdateInput struct {
	Name        *string
	Color       *string
	Description **string
}

type LabelService struct {
	repo  repositories.LabelRepository
	cache *cachepkg.Store
}

func NewLabelService(repo repositories.LabelRepository, cache *cachepkg.Store) *LabelService {
	return &LabelService{repo: repo, cache: cache}
}

func (s *LabelService) List(ctx context.Context, input LabelListInput) ([]LabelSummary, int64, *apperrors.AppError) {
	cacheKey := buildListCacheKey(
		"labels",
		intToCachePart(input.Page),
		intToCachePart(input.Limit),
		input.Search,
		input.SortBy,
		input.SortOrder,
	)
	if s.cache != nil {
		var cached struct {
			Items []LabelSummary `json:"items"`
			Total int64          `json:"total"`
		}
		if found, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && found {
			return cached.Items, cached.Total, nil
		}
	}

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

	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, struct {
			Items []LabelSummary `json:"items"`
			Total int64          `json:"total"`
		}{
			Items: items,
			Total: total,
		}, 10*time.Minute)
	}

	return items, total, nil
}

func (s *LabelService) Create(ctx context.Context, input LabelCreateInput) (*LabelSummary, *apperrors.AppError) {
	fields := apperrors.FieldErrors{}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		fields["name"] = "Label name is required."
	} else if utf8.RuneCountInString(name) > maxLabelNameLength {
		fields["name"] = "Label name must be 50 characters or fewer."
	}

	color := strings.TrimSpace(input.Color)
	if color == "" {
		fields["color"] = "Label color is required."
	} else if !labelColorPattern.MatchString(color) {
		fields["color"] = "Color must be a 6-digit hex code like #3B82F6."
	}

	description := normalizeLabelOptional(input.Description)
	if description != nil && utf8.RuneCountInString(*description) > maxLabelDescriptionLength {
		fields["description"] = "Description must be 1000 characters or fewer."
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", fields)
	}

	label := &models.Label{
		Name:        name,
		Color:       color,
		Description: description,
	}
	if err := s.repo.Create(ctx, label); err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Conflict("label name already exists", apperrors.FieldErrors{
				"name": "This label name is already in use.",
			})
		}
		return nil, apperrors.Internal("failed to create label")
	}

	s.invalidateLabelCaches(ctx)

	return &LabelSummary{
		ID:          label.ID,
		Name:        label.Name,
		Color:       label.Color,
		Description: label.Description,
		CreatedAt:   label.CreatedAt,
		UpdatedAt:   label.UpdatedAt,
	}, nil
}

func (s *LabelService) Get(ctx context.Context, id string) (*LabelDetail, *apperrors.AppError) {
	label, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("label not found")
		}
		return nil, apperrors.Internal("failed to load label")
	}

	usageCount, err := s.repo.UsageCountByID(ctx, id)
	if err != nil {
		return nil, apperrors.Internal("failed to load label usage")
	}

	return &LabelDetail{
		LabelSummary: LabelSummary{
			ID:          label.ID,
			Name:        label.Name,
			Color:       label.Color,
			Description: label.Description,
			CreatedAt:   label.CreatedAt,
			UpdatedAt:   label.UpdatedAt,
		},
		UsageCount: int(usageCount),
	}, nil
}

func (s *LabelService) Update(ctx context.Context, id string, input LabelUpdateInput) (*LabelSummary, *apperrors.AppError) {
	label, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("label not found")
		}
		return nil, apperrors.Internal("failed to load label")
	}

	fields := apperrors.FieldErrors{}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			fields["name"] = "Label name is required."
		} else if utf8.RuneCountInString(name) > maxLabelNameLength {
			fields["name"] = "Label name must be 50 characters or fewer."
		} else {
			label.Name = name
		}
	}
	if input.Color != nil {
		color := strings.TrimSpace(*input.Color)
		if color == "" {
			fields["color"] = "Label color is required."
		} else if !labelColorPattern.MatchString(color) {
			fields["color"] = "Color must be a 6-digit hex code like #3B82F6."
		} else {
			label.Color = color
		}
	}
	if input.Description != nil {
		description := normalizeLabelOptional(*input.Description)
		if description != nil && utf8.RuneCountInString(*description) > maxLabelDescriptionLength {
			fields["description"] = "Description must be 1000 characters or fewer."
		} else {
			label.Description = description
		}
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", fields)
	}

	if err := s.repo.Update(ctx, label); err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Conflict("label name already exists", apperrors.FieldErrors{
				"name": "This label name is already in use.",
			})
		}
		return nil, apperrors.Internal("failed to update label")
	}

	s.invalidateLabelCaches(ctx)

	return &LabelSummary{
		ID:          label.ID,
		Name:        label.Name,
		Color:       label.Color,
		Description: label.Description,
		CreatedAt:   label.CreatedAt,
		UpdatedAt:   label.UpdatedAt,
	}, nil
}

func (s *LabelService) Delete(ctx context.Context, id string) *apperrors.AppError {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return apperrors.NotFound("label not found")
		}
		return apperrors.Internal("failed to load label")
	}

	usageCount, err := s.repo.UsageCountByID(ctx, id)
	if err != nil {
		return apperrors.Internal("failed to validate label deletion")
	}
	if usageCount > 0 {
		return apperrors.Conflict("label cannot be deleted while in use", apperrors.FieldErrors{
			"id": "label is referenced by one or more issues",
		})
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return apperrors.NotFound("label not found")
		}
		return apperrors.Internal("failed to delete label")
	}
	s.invalidateLabelCaches(ctx)
	return nil
}

func (s *LabelService) invalidateLabelCaches(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.DeleteByPrefix(ctx, "labels:")
}

func normalizeLabelOptional(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}
