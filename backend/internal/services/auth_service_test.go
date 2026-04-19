package services

import (
	"context"
	"errors"
	"testing"

	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
)

type mockGoogleVerifier struct {
	verifyFn func(ctx context.Context, idToken, audience string) (*GoogleIdentity, error)
}

func (m *mockGoogleVerifier) VerifyIDToken(ctx context.Context, idToken, audience string) (*GoogleIdentity, error) {
	return m.verifyFn(ctx, idToken, audience)
}

type mockAuthRepo struct {
	createFn           func(ctx context.Context, user *models.User) error
	findByEmailFn      func(ctx context.Context, email string) (*models.User, error)
	findByGoogleSubFn  func(ctx context.Context, subject string) (*models.User, error)
	findByIDFn         func(ctx context.Context, id string) (*models.User, error)
	setGoogleSubjectFn func(ctx context.Context, userID, subject string) error
}

func (m *mockAuthRepo) Create(ctx context.Context, user *models.User) error {
	return m.createFn(ctx, user)
}

func (m *mockAuthRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.findByEmailFn(ctx, email)
}

func (m *mockAuthRepo) FindByGoogleSubject(ctx context.Context, subject string) (*models.User, error) {
	return m.findByGoogleSubFn(ctx, subject)
}

func (m *mockAuthRepo) SetGoogleSubject(ctx context.Context, userID, subject string) error {
	if m.setGoogleSubjectFn == nil {
		return nil
	}
	return m.setGoogleSubjectFn(ctx, userID, subject)
}

func (m *mockAuthRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	return m.findByIDFn(ctx, id)
}

func TestAuthService_LoginWithGoogle_LinksExistingEmail(t *testing.T) {
	t.Parallel()

	var linked string
	svc := NewAuthService(
		&mockAuthRepo{
			createFn: func(ctx context.Context, user *models.User) error { return nil },
			findByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					ID:    "61e4a2f3-0279-4ce7-82da-298492fdfdc0",
					Email: "alex@example.com",
					Name:  "Alex",
				}, nil
			},
			findByGoogleSubFn: func(ctx context.Context, subject string) (*models.User, error) {
				return nil, repositories.ErrNotFound
			},
			findByIDFn: func(ctx context.Context, id string) (*models.User, error) { return nil, repositories.ErrNotFound },
			setGoogleSubjectFn: func(ctx context.Context, userID, subject string) error {
				linked = subject
				return nil
			},
		},
		"0123456789abcdef0123456789abcdef",
		0,
		12,
		"google-client-id",
		&mockGoogleVerifier{
			verifyFn: func(ctx context.Context, idToken, audience string) (*GoogleIdentity, error) {
				return &GoogleIdentity{
					Subject: "google-sub-1",
					Email:   "alex@example.com",
					Name:    "Alex",
				}, nil
			},
		},
		nil,
	)

	session, appErr := svc.LoginWithGoogle(context.Background(), GoogleLoginInput{IDToken: "valid"})
	if appErr != nil {
		t.Fatalf("expected no app error, got %v", appErr)
	}
	if session == nil {
		t.Fatalf("expected session, got nil")
	}
	if linked != "google-sub-1" {
		t.Fatalf("expected google subject to be linked")
	}
}

func TestAuthService_LoginWithGoogle_CreatesUserWhenEmailMissing(t *testing.T) {
	t.Parallel()

	var created *models.User
	svc := NewAuthService(
		&mockAuthRepo{
			createFn: func(ctx context.Context, user *models.User) error {
				user.ID = "0a5f8f3d-3f86-4994-858d-74f1c0f40b12"
				created = user
				return nil
			},
			findByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
				return nil, repositories.ErrNotFound
			},
			findByGoogleSubFn: func(ctx context.Context, subject string) (*models.User, error) {
				return nil, repositories.ErrNotFound
			},
			findByIDFn: func(ctx context.Context, id string) (*models.User, error) { return nil, repositories.ErrNotFound },
		},
		"0123456789abcdef0123456789abcdef",
		0,
		12,
		"google-client-id",
		&mockGoogleVerifier{
			verifyFn: func(ctx context.Context, idToken, audience string) (*GoogleIdentity, error) {
				return &GoogleIdentity{
					Subject: "google-sub-2",
					Email:   "new@example.com",
					Name:    "New User",
				}, nil
			},
		},
		nil,
	)

	session, appErr := svc.LoginWithGoogle(context.Background(), GoogleLoginInput{IDToken: "valid"})
	if appErr != nil {
		t.Fatalf("expected no app error, got %v", appErr)
	}
	if session == nil || created == nil {
		t.Fatalf("expected created user and session")
	}
	if created.GoogleSubject == nil || *created.GoogleSubject != "google-sub-2" {
		t.Fatalf("expected created user to have google subject")
	}
}

func TestAuthService_LoginWithGoogle_InvalidToken(t *testing.T) {
	t.Parallel()

	svc := NewAuthService(
		&mockAuthRepo{
			createFn:          func(ctx context.Context, user *models.User) error { return nil },
			findByEmailFn:     func(ctx context.Context, email string) (*models.User, error) { return nil, repositories.ErrNotFound },
			findByGoogleSubFn: func(ctx context.Context, subject string) (*models.User, error) { return nil, repositories.ErrNotFound },
			findByIDFn:        func(ctx context.Context, id string) (*models.User, error) { return nil, repositories.ErrNotFound },
		},
		"0123456789abcdef0123456789abcdef",
		0,
		12,
		"google-client-id",
		&mockGoogleVerifier{
			verifyFn: func(ctx context.Context, idToken, audience string) (*GoogleIdentity, error) {
				return nil, errors.New("bad token")
			},
		},
		nil,
	)

	_, appErr := svc.LoginWithGoogle(context.Background(), GoogleLoginInput{IDToken: "invalid"})
	if appErr == nil {
		t.Fatalf("expected auth error")
	}
	if appErr.Code != "unauthorized" {
		t.Fatalf("expected unauthorized, got %s", appErr.Code)
	}
}
