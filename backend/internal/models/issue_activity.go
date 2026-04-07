package models

import "time"

const (
	IssueActivityCreated            = "created"
	IssueActivityUpdated            = "updated"
	IssueActivityTitleChanged       = "title_changed"
	IssueActivityDescriptionChanged = "description_changed"
	IssueActivityStatusChanged      = "status_changed"
	IssueActivityPriorityChanged    = "priority_changed"
	IssueActivityAssigneeChanged    = "assignee_changed"
	IssueActivitySprintChanged      = "sprint_changed"
	IssueActivityProjectChanged     = "project_changed"
	IssueActivityLabelAdded         = "label_added"
	IssueActivityLabelRemoved       = "label_removed"
	IssueActivityArchived           = "archived"
	IssueActivityRestored           = "restored"
)

// IssueActivity maps to the canonical issue_activities table.
type IssueActivity struct {
	ID        string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	IssueID   string    `gorm:"column:issue_id;type:uuid;not null"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null"`
	Action    string    `gorm:"column:action;type:varchar(50);not null"`
	FieldName *string   `gorm:"column:field_name;type:varchar(100)"`
	OldValue  *string   `gorm:"column:old_value;type:text"`
	NewValue  *string   `gorm:"column:new_value;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null"`
}

func (IssueActivity) TableName() string {
	return "issue_activities"
}
