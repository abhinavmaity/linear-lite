package models

import "time"

// IssueLabel maps to the canonical issue_labels join table.
type IssueLabel struct {
	IssueID   string    `gorm:"column:issue_id;type:uuid;primaryKey"`
	LabelID   string    `gorm:"column:label_id;type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null"`
}

func (IssueLabel) TableName() string {
	return "issue_labels"
}
