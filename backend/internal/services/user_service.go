package services

import (
	"context"
	"errors"

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
	repo repositories.UserReadRepository
}

func NewUserService(repo repositories.UserReadRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) List(ctx context.Context, input UserListInput) ([]UserSummary, int64, *apperrors.AppError) {
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
