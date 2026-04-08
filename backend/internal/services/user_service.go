package services

import (
	"context"
	"errors"
	"time"

	cachepkg "github.com/abhinavmaity/linear-lite/backend/internal/cache"
	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
)

type UserListInput struct {
	Page      int
	Limit     int
	Search    string
	SortBy    string
	SortOrder string
}

type UserService struct {
	repo  repositories.UserReadRepository
	cache *cachepkg.Store
}

func NewUserService(repo repositories.UserReadRepository, cache *cachepkg.Store) *UserService {
	return &UserService{repo: repo, cache: cache}
}

func (s *UserService) List(ctx context.Context, input UserListInput) ([]UserSummary, int64, *apperrors.AppError) {
	cacheKey := buildListCacheKey(
		"users",
		intToCachePart(input.Page),
		intToCachePart(input.Limit),
		input.Search,
		input.SortBy,
		input.SortOrder,
	)
	if s.cache != nil {
		var cached struct {
			Items []UserSummary `json:"items"`
			Total int64         `json:"total"`
		}
		if found, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && found {
			return cached.Items, cached.Total, nil
		}
	}

	users, total, err := s.repo.List(ctx, repositories.UserListFilter{
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
		return nil, 0, apperrors.Internal("failed to list users")
	}

	items := make([]UserSummary, 0, len(users))
	for _, user := range users {
		items = append(items, UserSummary{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, struct {
			Items []UserSummary `json:"items"`
			Total int64         `json:"total"`
		}{
			Items: items,
			Total: total,
		}, 5*time.Minute)
	}

	return items, total, nil
}

func (s *UserService) Get(ctx context.Context, id string) (*UserDetail, *apperrors.AppError) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("user not found")
		}
		return nil, apperrors.Internal("failed to load user")
	}

	stats, err := s.repo.IssueStatsByUserID(ctx, id)
	if err != nil {
		return nil, apperrors.Internal("failed to load user stats")
	}

	return &UserDetail{
		UserSummary: UserSummary{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Stats: UserStats{
			TotalCreated:       stats.TotalCreated,
			TotalAssigned:      stats.TotalAssigned,
			InProgressAssigned: stats.InProgressAssigned,
			DoneAssigned:       stats.DoneAssigned,
		},
	}, nil
}
