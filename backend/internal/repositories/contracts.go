package repositories

import (
	"context"

	"github.com/abhinavmaity/linear-lite/backend/internal/models"
)

// PaginationInput standardizes repository-level pagination options.
type PaginationInput struct {
	Page  int
	Limit int
}

// SortInput standardizes repository-level sort options.
type SortInput struct {
	By    string
	Order string
}

// UserListFilter contains filters for listing team users.
type UserListFilter struct {
	PaginationInput
	SortInput
	Search string
}

// ProjectListFilter contains filters for listing projects.
type ProjectListFilter struct {
	PaginationInput
	SortInput
	Search string
}

// SprintListFilter contains filters for listing sprints.
type SprintListFilter struct {
	PaginationInput
	SortInput
	Search    string
	ProjectID *string
	Status    *string
}

// LabelListFilter contains filters for listing labels.
type LabelListFilter struct {
	PaginationInput
	SortInput
	Search string
}

// IssueListFilter contains filters for list and board issue queries.
type IssueListFilter struct {
	PaginationInput
	SortInput
	Search          string
	Statuses        []string
	Priorities      []string
	AssigneeID      *string
	ProjectID       *string
	SprintID        *string
	LabelIDs        []string
	LabelMode       string
	IncludeArchived bool
}

// TxRunner abstracts transaction handling for service-layer operations.
type TxRunner interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// UserReadRepository defines read-only user queries used by selectors and stats.
type UserReadRepository interface {
	List(ctx context.Context, filter UserListFilter) ([]models.User, int64, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
}

// ProjectRepository defines project persistence contracts used in Milestone 3+.
type ProjectRepository interface {
	List(ctx context.Context, filter ProjectListFilter) ([]ProjectSummaryRow, int64, error)
	FindByID(ctx context.Context, id string) (*models.Project, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
}

// SprintRepository defines sprint persistence contracts used in Milestone 3+.
type SprintRepository interface {
	List(ctx context.Context, filter SprintListFilter) ([]SprintSummaryRow, int64, error)
	FindByID(ctx context.Context, id string) (*models.Sprint, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
}

// LabelRepository defines label persistence contracts used in Milestone 3+.
type LabelRepository interface {
	List(ctx context.Context, filter LabelListFilter) ([]models.Label, int64, error)
	FindByID(ctx context.Context, id string) (*models.Label, error)
	ExistsByIDs(ctx context.Context, ids []string) (bool, error)
}

// IssueRepository defines issue persistence contracts for list/detail/mutations.
type IssueRepository interface {
	List(ctx context.Context, filter IssueListFilter) ([]models.Issue, int64, error)
	FindByID(ctx context.Context, id string, includeArchived bool) (*models.Issue, error)
	Create(ctx context.Context, issue *models.Issue) error
	Update(ctx context.Context, issue *models.Issue) error
}

// IssueLabelRepository defines issue-label join-table operations.
type IssueLabelRepository interface {
	ListLabelIDsByIssueID(ctx context.Context, issueID string) ([]string, error)
	ReplaceIssueLabels(ctx context.Context, issueID string, labelIDs []string) error
}

// ActivityRepository defines issue activity persistence.
type ActivityRepository interface {
	ListByIssueID(ctx context.Context, issueID string) ([]models.IssueActivity, error)
	CreateMany(ctx context.Context, activities []models.IssueActivity) error
}
