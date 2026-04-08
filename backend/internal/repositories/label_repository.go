package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type LabelRepositoryDB struct {
	db *gorm.DB
}

func NewLabelRepository(db *gorm.DB) *LabelRepositoryDB {
	return &LabelRepositoryDB{db: db}
}

func (r *LabelRepositoryDB) FindByID(ctx context.Context, id string) (*models.Label, error) {
	var label models.Label
	query := r.db.WithContext(ctx).Where("id = ?", strings.TrimSpace(id)).First(&label)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if query.Error != nil {
		return nil, query.Error
	}
	return &label, nil
}

func (r *LabelRepositoryDB) ExistsByIDs(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	var count int64
	err := r.db.WithContext(ctx).Model(&models.Label{}).
		Where("id IN ?", ids).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return int(count) == len(ids), nil
}

func (r *LabelRepositoryDB) List(ctx context.Context, filter LabelListFilter) ([]models.Label, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Label{})

	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where("name ILIKE ?", pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []models.Label{}, 0, nil
	}

	sortBy := filter.By
	if sortBy == "" {
		sortBy = "name"
	}
	order := strings.ToLower(strings.TrimSpace(filter.Order))
	if order != "desc" {
		order = "asc"
	}

	orderClause := "name asc"
	switch sortBy {
	case "name":
		orderClause = "name " + order
	case "created_at":
		orderClause = "created_at " + order
	}

	var labels []models.Label
	err := query.Order(orderClause).
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Find(&labels).Error
	if err != nil {
		return nil, 0, err
	}

	return labels, total, nil
}

func (r *LabelRepositoryDB) FindByIDs(ctx context.Context, ids []string) ([]models.Label, error) {
	if len(ids) == 0 {
		return []models.Label{}, nil
	}

	var labels []models.Label
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&labels).Error
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (r *LabelRepositoryDB) Create(ctx context.Context, label *models.Label) error {
	if err := r.db.WithContext(ctx).Create(label).Error; err != nil {
		if isLabelUniqueViolation(err, "uq_labels_lower_name") {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *LabelRepositoryDB) Update(ctx context.Context, label *models.Label) error {
	if err := r.db.WithContext(ctx).Save(label).Error; err != nil {
		if isLabelUniqueViolation(err, "uq_labels_lower_name") {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *LabelRepositoryDB) Delete(ctx context.Context, id string) error {
	query := r.db.WithContext(ctx).Delete(&models.Label{}, "id = ?", strings.TrimSpace(id))
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *LabelRepositoryDB) UsageCountByID(ctx context.Context, id string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.IssueLabel{}).
		Where("label_id = ?", strings.TrimSpace(id)).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func isLabelUniqueViolation(err error, constraint string) bool {
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
