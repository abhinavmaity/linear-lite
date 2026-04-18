package services

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	cachepkg "github.com/abhinavmaity/linear-lite/backend/internal/cache"
	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 72
	maxEmailLength    = 255
	maxNameLength     = 255
)

type AuthUser struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthSession struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      AuthUser  `json:"user"`
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type UserAuthRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
}

type authClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	users      UserAuthRepository
	cache      *cachepkg.Store
	jwtSecret  []byte
	jwtTTL     time.Duration
	bcryptCost int
	now        func() time.Time
}

func NewAuthService(
	users UserAuthRepository,
	jwtSecret string,
	jwtTTL time.Duration,
	bcryptCost int,
	cache *cachepkg.Store,
) *AuthService {
	return &AuthService{
		users:      users,
		cache:      cache,
		jwtSecret:  []byte(strings.TrimSpace(jwtSecret)),
		jwtTTL:     jwtTTL,
		bcryptCost: bcryptCost,
		now:        time.Now,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthSession, *apperrors.AppError) {
	fields := apperrors.FieldErrors{}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		fields["name"] = "Name is required."
	} else if utf8.RuneCountInString(name) > maxNameLength {
		fields["name"] = "Name must be 255 characters or fewer."
	}

	emailInput := strings.TrimSpace(input.Email)
	if err := validateEmail(emailInput); err != nil {
		fields["email"] = err.Error()
	}

	if err := validatePassword(input.Password); err != nil {
		fields["password"] = err.Error()
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", fields)
	}

	email := strings.ToLower(emailInput)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.bcryptCost)
	if err != nil {
		return nil, apperrors.Internal("failed to process password")
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repositories.ErrEmailConflict) {
			return nil, apperrors.Conflict("email is already registered", apperrors.FieldErrors{
				"email": "This email address is already in use.",
			})
		}
		return nil, apperrors.Internal("failed to create user")
	}

	if s.cache != nil {
		_ = s.cache.DeleteByPrefix(ctx, "users:")
	}

	token, expiresAt, appErr := s.issueAccessToken(user.ID, user.Email)
	if appErr != nil {
		return nil, appErr
	}

	return &AuthSession{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      mapAuthUser(user),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthSession, *apperrors.AppError) {
	fields := apperrors.FieldErrors{}

	emailInput := strings.TrimSpace(input.Email)
	if err := validateEmail(emailInput); err != nil {
		fields["email"] = err.Error()
	}

	if err := validatePassword(input.Password); err != nil {
		fields["password"] = err.Error()
	}

	if len(fields) > 0 {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", fields)
	}

	email := strings.ToLower(emailInput)
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, invalidCredentialsError()
		}
		return nil, apperrors.Internal("failed to authenticate user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, invalidCredentialsError()
	}

	token, expiresAt, appErr := s.issueAccessToken(user.ID, user.Email)
	if appErr != nil {
		return nil, appErr
	}

	return &AuthSession{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      mapAuthUser(user),
	}, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*AuthUser, *apperrors.AppError) {
	cleanID := strings.TrimSpace(userID)
	if cleanID == "" {
		return nil, apperrors.Unauthorized("authentication is required")
	}
	if _, err := uuid.Parse(cleanID); err != nil {
		return nil, apperrors.Unauthorized("authentication is required")
	}

	user, err := s.users.FindByID(ctx, cleanID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("user not found")
		}
		return nil, apperrors.Internal("failed to load user")
	}

	summary := mapAuthUser(user)
	return &summary, nil
}

func (s *AuthService) issueAccessToken(userID, email string) (string, time.Time, *apperrors.AppError) {
	now := s.now().UTC()
	expiresAt := now.Add(s.jwtTTL)

	claims := authClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, apperrors.Internal("failed to issue access token")
	}

	return signedToken, expiresAt, nil
}

func mapAuthUser(user *models.User) AuthUser {
	return AuthUser{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("Email is required.")
	}
	if len(email) > maxEmailLength {
		return errors.New("Email must be 255 characters or fewer.")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return errors.New("Enter a valid email address.")
	}

	return nil
}

func validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < minPasswordLength || length > maxPasswordLength {
		return errors.New("Password must be 8-72 characters long.")
	}

	// bcrypt has a hard 72-byte input limit, so enforce it explicitly.
	if len([]byte(password)) > maxPasswordLength {
		return errors.New("Password must be 8-72 characters long.")
	}

	return nil
}

func invalidCredentialsError() *apperrors.AppError {
	return apperrors.Unauthorized("Email or password is incorrect.")
}
