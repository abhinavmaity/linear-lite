package services

import "time"

type IssueCounts struct {
	Total      int `json:"total"`
	Backlog    int `json:"backlog"`
	Todo       int `json:"todo"`
	InProgress int `json:"in_progress"`
	InReview   int `json:"in_review"`
	Done       int `json:"done"`
	Cancelled  int `json:"cancelled"`
}

type UserSummary struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStats struct {
	TotalCreated       int `json:"total_created"`
	TotalAssigned      int `json:"total_assigned"`
	InProgressAssigned int `json:"in_progress_assigned"`
	DoneAssigned       int `json:"done_assigned"`
}

type UserDetail struct {
	UserSummary
	Stats UserStats `json:"stats"`
}

type SprintSummary struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description *string     `json:"description"`
	ProjectID   string      `json:"project_id"`
	StartDate   string      `json:"start_date"`
	EndDate     string      `json:"end_date"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	IssueCounts IssueCounts `json:"issue_counts"`
}

type ProjectSummary struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  *string        `json:"description"`
	Key          string         `json:"key"`
	CreatedBy    string         `json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	IssueCounts  IssueCounts    `json:"issue_counts"`
	ActiveSprint *SprintSummary `json:"active_sprint"`
}

type ProjectDetail struct {
	ProjectSummary
	Creator UserSummary     `json:"creator"`
	Sprints []SprintSummary `json:"sprints"`
}

type SprintDetail struct {
	SprintSummary
	Project ProjectSummary `json:"project"`
}

type LabelSummary struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LabelDetail struct {
	LabelSummary
	UsageCount int `json:"usage_count"`
}

type IssueActivity struct {
	ID        string       `json:"id"`
	IssueID   string       `json:"issue_id"`
	UserID    string       `json:"user_id"`
	Action    string       `json:"action"`
	FieldName *string      `json:"field_name"`
	OldValue  *string      `json:"old_value"`
	NewValue  *string      `json:"new_value"`
	CreatedAt time.Time    `json:"created_at"`
	User      *UserSummary `json:"user"`
}

type IssueSummary struct {
	ID          string         `json:"id"`
	Identifier  string         `json:"identifier"`
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	Status      string         `json:"status"`
	Priority    string         `json:"priority"`
	ProjectID   string         `json:"project_id"`
	SprintID    *string        `json:"sprint_id"`
	AssigneeID  *string        `json:"assignee_id"`
	CreatedBy   string         `json:"created_by"`
	ArchivedAt  *time.Time     `json:"archived_at"`
	ArchivedBy  *string        `json:"archived_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Project     ProjectSummary `json:"project"`
	Sprint      *SprintSummary `json:"sprint"`
	Assignee    *UserSummary   `json:"assignee"`
	Creator     UserSummary    `json:"creator"`
	Labels      []LabelSummary `json:"labels"`
}

type IssueDetail struct {
	IssueSummary
	Activities []IssueActivity `json:"activities"`
}

type DashboardStats struct {
	TotalIssues    int             `json:"total_issues"`
	MyIssues       int             `json:"my_issues"`
	InProgress     int             `json:"in_progress"`
	DoneThisWeek   int             `json:"done_this_week"`
	ActiveSprint   *SprintSummary  `json:"active_sprint"`
	RecentActivity []IssueActivity `json:"recent_activity"`
}
