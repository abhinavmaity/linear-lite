package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if strings.TrimSpace(user.ID) == "" {
		user.ID = uuid.NewString()
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isUniqueViolation(err, "uq_users_lower_email") {
			return ErrEmailConflict
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := r.db.WithContext(ctx).
		Where("LOWER(email) = LOWER(?)", strings.TrimSpace(email)).
		First(&user)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if query.Error != nil {
		return nil, query.Error
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := r.db.WithContext(ctx).
		Where("id = ?", strings.TrimSpace(id)).
		First(&user)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if query.Error != nil {
		return nil, query.Error
	}

	return &user, nil
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code != "23505" {
			return false
		}
		if constraint == "" {
			return true
		}
		return pgErr.ConstraintName == constraint
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if string(pqErr.Code) != "23505" {
			return false
		}
		if constraint == "" {
			return true
		}
		return pqErr.Constraint == constraint
	}

	return false
}
