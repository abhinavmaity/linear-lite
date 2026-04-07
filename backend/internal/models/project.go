package models

import "time"

// Project maps to the canonical projects table managed by SQL migrations.
type Project struct {
	ID              string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Name            string    `gorm:"column:name;type:varchar(255);not null"`
	Description     *string   `gorm:"column:description;type:text"`
	Key             string    `gorm:"column:key;type:varchar(10);not null"`
	NextIssueNumber int       `gorm:"column:next_issue_number;type:integer;not null;default:1"`
	CreatedBy       string    `gorm:"column:created_by;type:uuid;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (Project) TableName() string {
	return "projects"
}
