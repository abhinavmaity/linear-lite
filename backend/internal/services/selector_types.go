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

type LabelSummary struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
