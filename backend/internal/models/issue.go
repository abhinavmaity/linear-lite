package models

import "time"

const (
	IssueStatusBacklog    = "backlog"
	IssueStatusTodo       = "todo"
	IssueStatusInProgress = "in_progress"
	IssueStatusInReview   = "in_review"
	IssueStatusDone       = "done"
	IssueStatusCancelled  = "cancelled"
)

const (
	IssuePriorityLow    = "low"
	IssuePriorityMedium = "medium"
	IssuePriorityHigh   = "high"
	IssuePriorityUrgent = "urgent"
)

// Issue maps to the canonical issues table managed by SQL migrations.
type Issue struct {
	ID           string     `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Identifier   string     `gorm:"column:identifier;type:varchar(32);not null"`
	Title        string     `gorm:"column:title;type:varchar(500);not null"`
	Description  *string    `gorm:"column:description;type:text"`
	Status       string     `gorm:"column:status;type:varchar(20);not null;default:backlog"`
	Priority     string     `gorm:"column:priority;type:varchar(10);not null;default:medium"`
	ProjectID    string     `gorm:"column:project_id;type:uuid;not null"`
	SprintID     *string    `gorm:"column:sprint_id;type:uuid"`
	AssigneeID   *string    `gorm:"column:assignee_id;type:uuid"`
	CreatedBy    string     `gorm:"column:created_by;type:uuid;not null"`
	ArchivedAt   *time.Time `gorm:"column:archived_at;type:timestamptz"`
	ArchivedBy   *string    `gorm:"column:archived_by;type:uuid"`
	SearchVector string     `gorm:"column:search_vector;type:tsvector;->"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (Issue) TableName() string {
	return "issues"
}
