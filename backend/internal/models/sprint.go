package models

import "time"

const (
	SprintStatusPlanned   = "planned"
	SprintStatusActive    = "active"
	SprintStatusCompleted = "completed"
)

// Sprint maps to the canonical sprints table managed by SQL migrations.
type Sprint struct {
	ID          string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	Description *string   `gorm:"column:description;type:text"`
	ProjectID   string    `gorm:"column:project_id;type:uuid;not null"`
	StartDate   time.Time `gorm:"column:start_date;type:date;not null"`
	EndDate     time.Time `gorm:"column:end_date;type:date;not null"`
	Status      string    `gorm:"column:status;type:varchar(20);not null;default:planned"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (Sprint) TableName() string {
	return "sprints"
}
