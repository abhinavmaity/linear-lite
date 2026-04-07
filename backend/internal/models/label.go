package models

import "time"

// Label maps to the canonical labels table managed by SQL migrations.
type Label struct {
	ID          string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(50);not null"`
	Color       string    `gorm:"column:color;type:varchar(7);not null"`
	Description *string   `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (Label) TableName() string {
	return "labels"
}
