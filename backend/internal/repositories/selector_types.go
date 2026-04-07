package repositories

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

type SprintSummaryRow struct {
	ID          string
	Name        string
	Description *string
	ProjectID   string
	StartDate   time.Time
	EndDate     time.Time
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IssueCounts IssueCounts
}

type ProjectSummaryRow struct {
	ID           string
	Name         string
	Description  *string
	Key          string
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IssueCounts  IssueCounts
	ActiveSprint *SprintSummaryRow
}
