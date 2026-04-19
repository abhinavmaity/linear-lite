package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/middleware"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type fakeAuthService struct {
	registerFn func(ctx context.Context, input services.RegisterInput) (*services.AuthSession, *apperrors.AppError)
	loginFn    func(ctx context.Context, input services.LoginInput) (*services.AuthSession, *apperrors.AppError)
	googleFn   func(ctx context.Context, input services.GoogleLoginInput) (*services.AuthSession, *apperrors.AppError)
	meFn       func(ctx context.Context, userID string) (*services.AuthUser, *apperrors.AppError)
}

func (f *fakeAuthService) Register(ctx context.Context, input services.RegisterInput) (*services.AuthSession, *apperrors.AppError) {
	return f.registerFn(ctx, input)
}

func (f *fakeAuthService) Login(ctx context.Context, input services.LoginInput) (*services.AuthSession, *apperrors.AppError) {
	return f.loginFn(ctx, input)
}

func (f *fakeAuthService) LoginWithGoogle(ctx context.Context, input services.GoogleLoginInput) (*services.AuthSession, *apperrors.AppError) {
	if f.googleFn == nil {
		return nil, nil
	}
	return f.googleFn(ctx, input)
}

func (f *fakeAuthService) Me(ctx context.Context, userID string) (*services.AuthUser, *apperrors.AppError) {
	return f.meFn(ctx, userID)
}

type fakeProjectService struct {
	createFn func(ctx context.Context, actorID string, input services.ProjectCreateInput) (*services.ProjectDetail, *apperrors.AppError)
}

func (f *fakeProjectService) List(ctx context.Context, input services.ProjectListInput) ([]services.ProjectSummary, int64, *apperrors.AppError) {
	return nil, 0, nil
}

func (f *fakeProjectService) Create(ctx context.Context, actorID string, input services.ProjectCreateInput) (*services.ProjectDetail, *apperrors.AppError) {
	return f.createFn(ctx, actorID, input)
}

func (f *fakeProjectService) Get(ctx context.Context, id string) (*services.ProjectDetail, *apperrors.AppError) {
	return nil, nil
}

func (f *fakeProjectService) Update(ctx context.Context, id string, input services.ProjectUpdateInput) (*services.ProjectDetail, *apperrors.AppError) {
	return nil, nil
}

func (f *fakeProjectService) Delete(ctx context.Context, id string) *apperrors.AppError {
	return nil
}

type fakeUserService struct {
	listFn func(ctx context.Context, input services.UserListInput) ([]services.UserSummary, int64, *apperrors.AppError)
	getFn  func(ctx context.Context, id string) (*services.UserDetail, *apperrors.AppError)
}

func (f *fakeUserService) List(ctx context.Context, input services.UserListInput) ([]services.UserSummary, int64, *apperrors.AppError) {
	return f.listFn(ctx, input)
}

func (f *fakeUserService) Get(ctx context.Context, id string) (*services.UserDetail, *apperrors.AppError) {
	if f.getFn == nil {
		return nil, nil
	}
	return f.getFn(ctx, id)
}

type fakeDashboardService struct {
	getStatsFn func(ctx context.Context, userID string) (*services.DashboardStats, *apperrors.AppError)
}

func (f *fakeDashboardService) GetStats(ctx context.Context, userID string) (*services.DashboardStats, *apperrors.AppError) {
	return f.getStatsFn(ctx, userID)
}

type fakeIssueService struct {
	listFn    func(ctx context.Context, input services.IssueListInput) ([]services.IssueSummary, int64, *apperrors.AppError)
	getFn     func(ctx context.Context, id string, includeArchived bool) (*services.IssueDetail, *apperrors.AppError)
	createFn  func(ctx context.Context, actorID string, input services.CreateIssueInput) (*services.IssueDetail, *apperrors.AppError)
	updateFn  func(ctx context.Context, actorID string, input services.UpdateIssueInput) (*services.IssueDetail, *apperrors.AppError)
	archiveFn func(ctx context.Context, actorID string, id string) *apperrors.AppError
}

func (f *fakeIssueService) List(ctx context.Context, input services.IssueListInput) ([]services.IssueSummary, int64, *apperrors.AppError) {
	if f.listFn == nil {
		return nil, 0, nil
	}
	return f.listFn(ctx, input)
}

func (f *fakeIssueService) Get(ctx context.Context, id string, includeArchived bool) (*services.IssueDetail, *apperrors.AppError) {
	return f.getFn(ctx, id, includeArchived)
}

func (f *fakeIssueService) Create(ctx context.Context, actorID string, input services.CreateIssueInput) (*services.IssueDetail, *apperrors.AppError) {
	if f.createFn == nil {
		return nil, nil
	}
	return f.createFn(ctx, actorID, input)
}

func (f *fakeIssueService) Update(ctx context.Context, actorID string, input services.UpdateIssueInput) (*services.IssueDetail, *apperrors.AppError) {
	if f.updateFn == nil {
		return nil, nil
	}
	return f.updateFn(ctx, actorID, input)
}

func (f *fakeIssueService) Archive(ctx context.Context, actorID string, id string) *apperrors.AppError {
	if f.archiveFn == nil {
		return nil
	}
	return f.archiveFn(ctx, actorID, id)
}

func TestAuthRegister_InvalidJSON_ReturnsValidationEnvelope(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&fakeAuthService{
		registerFn: func(ctx context.Context, input services.RegisterInput) (*services.AuthSession, *apperrors.AppError) {
			t.Fatalf("register service should not be called for invalid json")
			return nil, nil
		},
		loginFn: func(ctx context.Context, input services.LoginInput) (*services.AuthSession, *apperrors.AppError) {
			return nil, nil
		},
		meFn: func(ctx context.Context, userID string) (*services.AuthUser, *apperrors.AppError) {
			return nil, nil
		},
	})

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"name":"bad"`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}
	errObj := payload["error"].(map[string]any)
	if got := errObj["code"]; got != "validation_error" {
		t.Fatalf("expected error.code validation_error, got %v", got)
	}
}

func TestAuthRegister_ConflictFromService_PropagatesStatusAndCode(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&fakeAuthService{
		registerFn: func(ctx context.Context, input services.RegisterInput) (*services.AuthSession, *apperrors.AppError) {
			return nil, apperrors.Conflict("email already exists", apperrors.FieldErrors{"email": "already in use"})
		},
		loginFn: func(ctx context.Context, input services.LoginInput) (*services.AuthSession, *apperrors.AppError) {
			return nil, nil
		},
		meFn: func(ctx context.Context, userID string) (*services.AuthUser, *apperrors.AppError) {
			return nil, nil
		},
	})

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"name":"A","email":"a@example.com","password":"Password123"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, resp.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}
	errObj := payload["error"].(map[string]any)
	if got := errObj["code"]; got != "conflict" {
		t.Fatalf("expected error.code conflict, got %v", got)
	}
}

func TestProjectCreate_ValidationErrorContainsRequestID(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	handler := NewProjectHandler(&fakeProjectService{
		createFn: func(ctx context.Context, actorID string, input services.ProjectCreateInput) (*services.ProjectDetail, *apperrors.AppError) {
			return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{
				"key": "must match ^[A-Z0-9]{2,10}$",
			})
		},
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyRequestID, "req-test-1")
		c.Next()
	})
	router.POST("/projects", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBufferString(`{"name":"X","key":"bad-key"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}
	errObj := payload["error"].(map[string]any)
	if got := errObj["code"]; got != "validation_error" {
		t.Fatalf("expected error.code validation_error, got %v", got)
	}
	if got := errObj["request_id"]; got != "req-test-1" {
		t.Fatalf("expected request_id req-test-1, got %v", got)
	}
}

func TestUserList_EmptyItemsSerializesAsJSONList(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	handler := NewUserHandler(&fakeUserService{
		listFn: func(ctx context.Context, input services.UserListInput) ([]services.UserSummary, int64, *apperrors.AppError) {
			return []services.UserSummary{}, 0, nil
		},
	})

	router := gin.New()
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	items, ok := payload["items"].([]any)
	if !ok {
		t.Fatalf("expected items to be a json array, got %T", payload["items"])
	}
	if len(items) != 0 {
		t.Fatalf("expected empty items array, got %d entries", len(items))
	}
}

func TestDashboardStats_EmptyRecentActivitySerializesAsJSONList(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	handler := NewDashboardHandler(&fakeDashboardService{
		getStatsFn: func(ctx context.Context, userID string) (*services.DashboardStats, *apperrors.AppError) {
			return &services.DashboardStats{
				RecentActivity: []services.IssueActivity{},
			}, nil
		},
	})

	router := gin.New()
	router.GET("/dashboard/stats", handler.Stats)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/stats", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", payload["data"])
	}
	activity, ok := data["recent_activity"].([]any)
	if !ok {
		t.Fatalf("expected recent_activity to be a json array, got %T", data["recent_activity"])
	}
	if len(activity) != 0 {
		t.Fatalf("expected empty recent_activity array, got %d entries", len(activity))
	}
}

func TestIssueGet_EmptyLabelsAndActivitiesSerializeAsJSONLists(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	handler := NewIssueHandler(&fakeIssueService{
		getFn: func(ctx context.Context, id string, includeArchived bool) (*services.IssueDetail, *apperrors.AppError) {
			return &services.IssueDetail{
				IssueSummary: services.IssueSummary{
					ID:         id,
					Identifier: "PRJ-1",
					Labels:     []services.LabelSummary{},
				},
				Activities: []services.IssueActivity{},
			}, nil
		},
	})

	router := gin.New()
	router.GET("/issues/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/issues/00000000-0000-0000-0000-000000000001", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", payload["data"])
	}

	labels, ok := data["labels"].([]any)
	if !ok {
		t.Fatalf("expected labels to be a json array, got %T", data["labels"])
	}
	if len(labels) != 0 {
		t.Fatalf("expected empty labels array, got %d entries", len(labels))
	}

	activities, ok := data["activities"].([]any)
	if !ok {
		t.Fatalf("expected activities to be a json array, got %T", data["activities"])
	}
	if len(activities) != 0 {
		t.Fatalf("expected empty activities array, got %d entries", len(activities))
	}
}
