package services

import (
	"context"
	"strings"
	"time"

	cachepkg "github.com/abhinavmaity/linear-lite/backend/internal/cache"
	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
)

type DashboardService struct {
	issues  *repositories.IssueRepositoryDB
	sprints repositories.SprintRepository
	users   repositories.UserReadRepository
	cache   *cachepkg.Store
	now     func() time.Time
}

func NewDashboardService(
	issues *repositories.IssueRepositoryDB,
	sprints repositories.SprintRepository,
	users repositories.UserReadRepository,
	cache *cachepkg.Store,
) *DashboardService {
	return &DashboardService{
		issues:  issues,
		sprints: sprints,
		users:   users,
		cache:   cache,
		now:     time.Now,
	}
}

func (s *DashboardService) GetStats(ctx context.Context, userID string) (*DashboardStats, *apperrors.AppError) {
	cleanUserID := strings.TrimSpace(userID)
	if cleanUserID == "" {
		return nil, apperrors.Unauthorized("authentication is required")
	}
	cacheKey := "dashboard:stats:" + cleanUserID
	if s.cache != nil {
		var cached DashboardStats
		if found, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && found {
			return &cached, nil
		}
	}

	doneSince := s.now().UTC().AddDate(0, 0, -7)
	metrics, err := s.issues.DashboardMetrics(ctx, cleanUserID, doneSince)
	if err != nil {
		return nil, apperrors.Internal("failed to load dashboard metrics")
	}

	var activeSprint *SprintSummary
	activeSprintID, err := s.issues.DashboardActiveSprintID(ctx)
	if err != nil {
		return nil, apperrors.Internal("failed to load active sprint")
	}
	if activeSprintID != nil {
		sprintMap, sprintErr := s.sprints.SummariesByIDs(ctx, []string{*activeSprintID})
		if sprintErr != nil {
			return nil, apperrors.Internal("failed to load active sprint")
		}
		if sprintRow, ok := sprintMap[*activeSprintID]; ok {
			activeSprint = mapSprintRow(sprintRow)
		}
	}

	activityRows, err := s.issues.ListRecentActivitiesForDashboard(ctx, 10)
	if err != nil {
		return nil, apperrors.Internal("failed to load recent activity")
	}

	recentActivity, appErr := s.hydrateDashboardActivities(ctx, activityRows)
	if appErr != nil {
		return nil, appErr
	}

	result := &DashboardStats{
		TotalIssues:    metrics.TotalIssues,
		MyIssues:       metrics.MyIssues,
		InProgress:     metrics.InProgress,
		DoneThisWeek:   metrics.DoneThisWeek,
		ActiveSprint:   activeSprint,
		RecentActivity: recentActivity,
	}

	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, cacheKey, result, 30*time.Second)
	}

	return result, nil
}

func (s *DashboardService) hydrateDashboardActivities(
	ctx context.Context,
	rows []models.IssueActivity,
) ([]IssueActivity, *apperrors.AppError) {
	if len(rows) == 0 {
		return []IssueActivity{}, nil
	}

	userSet := make(map[string]struct{}, len(rows))
	userIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		if _, ok := userSet[row.UserID]; ok {
			continue
		}
		userSet[row.UserID] = struct{}{}
		userIDs = append(userIDs, row.UserID)
	}

	users, err := s.users.FindByIDs(ctx, userIDs)
	if err != nil {
		return nil, apperrors.Internal("failed to load activity users")
	}
	userMap := make(map[string]UserSummary, len(users))
	for _, user := range users {
		userMap[user.ID] = UserSummary{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	items := make([]IssueActivity, 0, len(rows))
	for _, row := range rows {
		var user *UserSummary
		if u, ok := userMap[row.UserID]; ok {
			copy := u
			user = &copy
		}
		items = append(items, IssueActivity{
			ID:        row.ID,
			IssueID:   row.IssueID,
			UserID:    row.UserID,
			Action:    row.Action,
			FieldName: row.FieldName,
			OldValue:  row.OldValue,
			NewValue:  row.NewValue,
			CreatedAt: row.CreatedAt,
			User:      user,
		})
	}
	return items, nil
}
