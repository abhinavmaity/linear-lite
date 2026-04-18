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
	meFn       func(ctx context.Context, userID string) (*services.AuthUser, *apperrors.AppError)
}

func (f *fakeAuthService) Register(ctx context.Context, input services.RegisterInput) (*services.AuthSession, *apperrors.AppError) {
	return f.registerFn(ctx, input)
}

func (f *fakeAuthService) Login(ctx context.Context, input services.LoginInput) (*services.AuthSession, *apperrors.AppError) {
	return f.loginFn(ctx, input)
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
