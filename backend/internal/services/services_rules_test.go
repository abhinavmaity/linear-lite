package services

import (
	"context"
	"testing"
	"time"

	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
)

type mockProjectRepo struct {
	findByIDFn             func(ctx context.Context, id string) (*models.Project, error)
	countIssuesByProjectFn func(ctx context.Context, id string) (int64, error)
	countSprintsByProjFn   func(ctx context.Context, id string) (int64, error)
	deleteFn               func(ctx context.Context, id string) error
	existsByIDFn           func(ctx context.Context, id string) (bool, error)
}

func (m *mockProjectRepo) List(ctx context.Context, filter repositories.ProjectListFilter) ([]repositories.ProjectSummaryRow, int64, error) {
	return nil, 0, nil
}
func (m *mockProjectRepo) FindByID(ctx context.Context, id string) (*models.Project, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockProjectRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	if m.existsByIDFn == nil {
		return false, nil
	}
	return m.existsByIDFn(ctx, id)
}
func (m *mockProjectRepo) SummariesByIDs(ctx context.Context, ids []string) (map[string]repositories.ProjectSummaryRow, error) {
	return map[string]repositories.ProjectSummaryRow{}, nil
}
func (m *mockProjectRepo) ListSprintsByProjectID(ctx context.Context, projectID string) ([]repositories.SprintSummaryRow, error) {
	return nil, nil
}
func (m *mockProjectRepo) Create(ctx context.Context, project *models.Project) error { return nil }
func (m *mockProjectRepo) Update(ctx context.Context, project *models.Project) error { return nil }
func (m *mockProjectRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn == nil {
		return nil
	}
	return m.deleteFn(ctx, id)
}
func (m *mockProjectRepo) CountIssuesByProjectID(ctx context.Context, id string) (int64, error) {
	return m.countIssuesByProjectFn(ctx, id)
}
func (m *mockProjectRepo) CountSprintsByProjectID(ctx context.Context, id string) (int64, error) {
	if m.countSprintsByProjFn == nil {
		return 0, nil
	}
	return m.countSprintsByProjFn(ctx, id)
}

type mockUserReadRepo struct{}

func (m *mockUserReadRepo) List(ctx context.Context, filter repositories.UserListFilter) ([]models.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserReadRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	return nil, repositories.ErrNotFound
}
func (m *mockUserReadRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	return false, nil
}
func (m *mockUserReadRepo) FindByIDs(ctx context.Context, ids []string) ([]models.User, error) {
	return nil, nil
}
func (m *mockUserReadRepo) IssueStatsByUserID(ctx context.Context, id string) (repositories.UserIssueStats, error) {
	return repositories.UserIssueStats{}, nil
}

type mockSprintRepo struct {
	findByIDFn       func(ctx context.Context, id string) (*models.Sprint, error)
	createFn         func(ctx context.Context, sprint *models.Sprint) error
	countIssuesByFn  func(ctx context.Context, id string) (int64, error)
	deleteFn         func(ctx context.Context, id string) error
	listFn           func(ctx context.Context, filter repositories.SprintListFilter) ([]repositories.SprintSummaryRow, int64, error)
	summariesByIDsFn func(ctx context.Context, ids []string) (map[string]repositories.SprintSummaryRow, error)
}

func (m *mockSprintRepo) List(ctx context.Context, filter repositories.SprintListFilter) ([]repositories.SprintSummaryRow, int64, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filter)
	}
	return nil, 0, nil
}
func (m *mockSprintRepo) FindByID(ctx context.Context, id string) (*models.Sprint, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockSprintRepo) ExistsByID(ctx context.Context, id string) (bool, error) { return false, nil }
func (m *mockSprintRepo) SummariesByIDs(ctx context.Context, ids []string) (map[string]repositories.SprintSummaryRow, error) {
	if m.summariesByIDsFn != nil {
		return m.summariesByIDsFn(ctx, ids)
	}
	return map[string]repositories.SprintSummaryRow{}, nil
}
func (m *mockSprintRepo) Create(ctx context.Context, sprint *models.Sprint) error {
	if m.createFn == nil {
		return nil
	}
	return m.createFn(ctx, sprint)
}
func (m *mockSprintRepo) Update(ctx context.Context, sprint *models.Sprint) error { return nil }
func (m *mockSprintRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn == nil {
		return nil
	}
	return m.deleteFn(ctx, id)
}
func (m *mockSprintRepo) CountIssuesBySprintID(ctx context.Context, id string) (int64, error) {
	if m.countIssuesByFn == nil {
		return 0, nil
	}
	return m.countIssuesByFn(ctx, id)
}

func TestProjectService_Delete_BlocksWhenIssuesExist(t *testing.T) {
	t.Parallel()

	svc := NewProjectService(
		&mockProjectRepo{
			findByIDFn: func(ctx context.Context, id string) (*models.Project, error) {
				return &models.Project{ID: id, Name: "P", Key: "PRJ"}, nil
			},
			countIssuesByProjectFn: func(ctx context.Context, id string) (int64, error) {
				return 3, nil
			},
			countSprintsByProjFn: func(ctx context.Context, id string) (int64, error) {
				return 0, nil
			},
		},
		&mockUserReadRepo{},
		nil,
	)

	appErr := svc.Delete(context.Background(), "project-1")
	if appErr == nil {
		t.Fatalf("expected conflict error")
	}
	if appErr.Code != "conflict" {
		t.Fatalf("expected conflict code, got %s", appErr.Code)
	}
}

func TestProjectService_Update_BlocksKeyChangeWhenIssuesExist(t *testing.T) {
	t.Parallel()

	svc := NewProjectService(
		&mockProjectRepo{
			findByIDFn: func(ctx context.Context, id string) (*models.Project, error) {
				return &models.Project{ID: id, Name: "P", Key: "PRJ"}, nil
			},
			countIssuesByProjectFn: func(ctx context.Context, id string) (int64, error) {
				return 1, nil
			},
		},
		&mockUserReadRepo{},
		nil,
	)

	nextKey := "NEXT"
	_, appErr := svc.Update(context.Background(), "project-1", ProjectUpdateInput{
		Key: &nextKey,
	})
	if appErr == nil {
		t.Fatalf("expected conflict error")
	}
	if appErr.Code != "conflict" {
		t.Fatalf("expected conflict code, got %s", appErr.Code)
	}
}

func TestSprintService_Create_MapsActiveConflict(t *testing.T) {
	t.Parallel()

	svc := NewSprintService(
		&mockSprintRepo{
			createFn: func(ctx context.Context, sprint *models.Sprint) error {
				return repositories.ErrConflict
			},
		},
		&mockProjectRepo{
			existsByIDFn: func(ctx context.Context, id string) (bool, error) {
				return true, nil
			},
		},
		nil,
	)

	active := models.SprintStatusActive
	_, appErr := svc.Create(context.Background(), SprintCreateInput{
		Name:      "Sprint 1",
		ProjectID: "9c1e4f4a-bcda-42db-89e1-3a7c7ed3809c",
		StartDate: "2026-04-12",
		EndDate:   "2026-04-19",
		Status:    &active,
	})
	if appErr == nil {
		t.Fatalf("expected conflict error")
	}
	if appErr.Code != "conflict" {
		t.Fatalf("expected conflict code, got %s", appErr.Code)
	}
}

func TestSprintService_Delete_BlocksActiveAndDependentSprint(t *testing.T) {
	t.Parallel()

	svc := NewSprintService(
		&mockSprintRepo{
			findByIDFn: func(ctx context.Context, id string) (*models.Sprint, error) {
				now := time.Now().UTC()
				return &models.Sprint{
					ID:        id,
					Name:      "Active Sprint",
					ProjectID: "9c1e4f4a-bcda-42db-89e1-3a7c7ed3809c",
					StartDate: now,
					EndDate:   now.Add(24 * time.Hour),
					Status:    models.SprintStatusActive,
				}, nil
			},
		},
		&mockProjectRepo{},
		nil,
	)

	appErr := svc.Delete(context.Background(), "sprint-1")
	if appErr == nil {
		t.Fatalf("expected conflict error for active sprint")
	}
	if appErr.Code != "conflict" {
		t.Fatalf("expected conflict code, got %s", appErr.Code)
	}
}

func TestSprintService_Delete_BlocksWhenIssuesReferenceSprint(t *testing.T) {
	t.Parallel()

	svc := NewSprintService(
		&mockSprintRepo{
			findByIDFn: func(ctx context.Context, id string) (*models.Sprint, error) {
				now := time.Now().UTC()
				return &models.Sprint{
					ID:        id,
					Name:      "Planned Sprint",
					ProjectID: "9c1e4f4a-bcda-42db-89e1-3a7c7ed3809c",
					StartDate: now,
					EndDate:   now.Add(24 * time.Hour),
					Status:    models.SprintStatusPlanned,
				}, nil
			},
			countIssuesByFn: func(ctx context.Context, id string) (int64, error) {
				return 2, nil
			},
		},
		&mockProjectRepo{},
		nil,
	)

	appErr := svc.Delete(context.Background(), "sprint-1")
	if appErr == nil {
		t.Fatalf("expected conflict error for dependent issues")
	}
	if appErr.Code != "conflict" {
		t.Fatalf("expected conflict code, got %s", appErr.Code)
	}
}
