package services

import (
	"context"
	"errors"
	"strings"

	cachepkg "github.com/abhinavmaity/linear-lite/backend/internal/cache"
	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/models"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
)

const (
	maxIssueTitleLength       = 500
	maxIssueDescriptionLength = 50000
)

type IssueListInput struct {
	Page            int
	Limit           int
	SortBy          string
	SortOrder       string
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

type CreateIssueInput struct {
	Title       string
	Description *string
	Status      *string
	Priority    *string
	ProjectID   string
	SprintID    *string
	AssigneeID  *string
	LabelIDs    []string
}

type UpdateIssueInput struct {
	ID          string
	Title       *string
	Description **string
	Status      *string
	Priority    *string
	ProjectID   *string
	SprintID    **string
	AssigneeID  **string
	LabelIDs    *[]string
	Archived    *bool
}

type IssueService struct {
	issues   *repositories.IssueRepositoryDB
	users    *repositories.UserRepository
	projects *repositories.ProjectRepositoryDB
	sprints  *repositories.SprintRepositoryDB
	labels   *repositories.LabelRepositoryDB
	cache    *cachepkg.Store
}

func NewIssueService(
	issues *repositories.IssueRepositoryDB,
	users *repositories.UserRepository,
	projects *repositories.ProjectRepositoryDB,
	sprints *repositories.SprintRepositoryDB,
	labels *repositories.LabelRepositoryDB,
	cache *cachepkg.Store,
) *IssueService {
	return &IssueService{
		issues:   issues,
		users:    users,
		projects: projects,
		sprints:  sprints,
		labels:   labels,
		cache:    cache,
	}
}

func (s *IssueService) List(ctx context.Context, input IssueListInput) ([]IssueSummary, int64, *apperrors.AppError) {
	rows, total, err := s.issues.List(ctx, repositories.IssueListFilter{
		PaginationInput: repositories.PaginationInput{
			Page:  input.Page,
			Limit: input.Limit,
		},
		SortInput: repositories.SortInput{
			By:    input.SortBy,
			Order: input.SortOrder,
		},
		Search:          input.Search,
		Statuses:        input.Statuses,
		Priorities:      input.Priorities,
		AssigneeID:      input.AssigneeID,
		ProjectID:       input.ProjectID,
		SprintID:        input.SprintID,
		LabelIDs:        input.LabelIDs,
		LabelMode:       input.LabelMode,
		IncludeArchived: input.IncludeArchived,
	})
	if err != nil {
		return nil, 0, apperrors.Internal("failed to list issues")
	}

	items, hydrateErr := s.hydrateIssues(ctx, rows)
	if hydrateErr != nil {
		return nil, 0, hydrateErr
	}
	return items, total, nil
}

func (s *IssueService) Get(ctx context.Context, id string, includeArchived bool) (*IssueDetail, *apperrors.AppError) {
	row, err := s.issues.FindByID(ctx, id, includeArchived)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("issue not found")
		}
		return nil, apperrors.Internal("failed to load issue")
	}

	summaries, hydrateErr := s.hydrateIssues(ctx, []models.Issue{*row})
	if hydrateErr != nil {
		return nil, hydrateErr
	}
	if len(summaries) == 0 {
		return nil, apperrors.NotFound("issue not found")
	}

	activities, err := s.hydrateIssueActivities(ctx, row.ID)
	if err != nil {
		return nil, apperrors.Internal("failed to load issue activities")
	}

	return &IssueDetail{
		IssueSummary: summaries[0],
		Activities:   activities,
	}, nil
}

func (s *IssueService) Create(ctx context.Context, actorID string, input CreateIssueInput) (*IssueDetail, *apperrors.AppError) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{"title": "Title is required."})
	}
	if len(title) > maxIssueTitleLength {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{"title": "Title must be 500 characters or fewer."})
	}

	description := normalizeOptional(input.Description)
	if description != nil && len(*description) > maxIssueDescriptionLength {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{"description": "Description must be 50000 characters or fewer."})
	}

	status := models.IssueStatusBacklog
	if input.Status != nil {
		status = *input.Status
	}
	priority := models.IssuePriorityMedium
	if input.Priority != nil {
		priority = *input.Priority
	}

	if err := s.validateIssueReferences(ctx, input.ProjectID, input.SprintID, input.AssigneeID, input.LabelIDs); err != nil {
		return nil, err
	}

	created, err := s.issues.CreateWithRelations(ctx, repositories.CreateIssueInput{
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		ProjectID:   input.ProjectID,
		SprintID:    input.SprintID,
		AssigneeID:  input.AssigneeID,
		CreatedBy:   actorID,
		LabelIDs:    input.LabelIDs,
	})
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("project not found")
		}
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Conflict("could not generate issue identifier", nil)
		}
		return nil, apperrors.Internal("failed to create issue")
	}

	s.invalidateIssueMutationCaches(ctx)

	return s.Get(ctx, created.ID, true)
}

func (s *IssueService) Update(ctx context.Context, actorID string, input UpdateIssueInput) (*IssueDetail, *apperrors.AppError) {
	if input.Archived != nil && *input.Archived {
		return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{
			"archived": "Use the archive action instead of setting archived=true.",
		})
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{"title": "Title is required."})
		}
		if len(title) > maxIssueTitleLength {
			return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{"title": "Title must be 500 characters or fewer."})
		}
		*input.Title = title
	}
	if input.Description != nil {
		desc := normalizeOptional(*input.Description)
		if desc != nil && len(*desc) > maxIssueDescriptionLength {
			return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{"description": "Description must be 50000 characters or fewer."})
		}
		*input.Description = desc
	}

	var projectIDForValidation *string
	if input.ProjectID != nil {
		projectIDForValidation = input.ProjectID
	}
	if projectIDForValidation != nil || input.SprintID != nil || input.AssigneeID != nil || input.LabelIDs != nil {
		current, err := s.issues.FindByID(ctx, input.ID, true)
		if err != nil {
			if errors.Is(err, repositories.ErrNotFound) {
				return nil, apperrors.NotFound("issue not found")
			}
			return nil, apperrors.Internal("failed to load issue")
		}

		targetProjectID := current.ProjectID
		if input.ProjectID != nil {
			targetProjectID = *input.ProjectID
		}
		var sprintID *string
		if input.SprintID != nil {
			sprintID = *input.SprintID
		} else {
			sprintID = current.SprintID
		}
		var assigneeID *string
		if input.AssigneeID != nil {
			assigneeID = *input.AssigneeID
		} else {
			assigneeID = current.AssigneeID
		}
		labelIDs := []string{}
		if input.LabelIDs != nil {
			labelIDs = *input.LabelIDs
		}

		if err := s.validateIssueReferences(ctx, targetProjectID, sprintID, assigneeID, labelIDs); err != nil {
			return nil, err
		}
		if input.ProjectID != nil && input.SprintID == nil && current.SprintID != nil {
			sprint, err := s.sprints.FindByID(ctx, *current.SprintID)
			if err != nil {
				return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{
					"sprint_id": "Select a sprint that belongs to the selected project.",
				})
			}
			if sprint.ProjectID != targetProjectID {
				return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{
					"sprint_id": "Clear or replace sprint when changing the project.",
				})
			}
		}
	}

	restore := input.Archived != nil && !*input.Archived
	updated, _, err := s.issues.UpdateWithRelations(ctx, repositories.UpdateIssueInput{
		ID:          input.ID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		ProjectID:   input.ProjectID,
		SprintID:    input.SprintID,
		AssigneeID:  input.AssigneeID,
		ActorID:     actorID,
		LabelIDs:    input.LabelIDs,
		Restore:     restore,
	})
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, apperrors.NotFound("issue not found")
		}
		if errors.Is(err, repositories.ErrConflict) {
			return nil, apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{
				"archived": "Only archived issues can be restored.",
			})
		}
		return nil, apperrors.Internal("failed to update issue")
	}

	s.invalidateIssueMutationCaches(ctx)

	return s.Get(ctx, updated.ID, true)
}

func (s *IssueService) Archive(ctx context.Context, actorID string, id string) *apperrors.AppError {
	_, err := s.issues.Archive(ctx, id, actorID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return apperrors.NotFound("issue not found")
		}
		return apperrors.Internal("failed to archive issue")
	}
	s.invalidateIssueMutationCaches(ctx)
	return nil
}

func (s *IssueService) invalidateIssueMutationCaches(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.DeleteByPrefix(ctx, "dashboard:")
	_ = s.cache.DeleteByPrefix(ctx, "projects:")
	_ = s.cache.DeleteByPrefix(ctx, "sprints:")
}

func (s *IssueService) validateIssueReferences(
	ctx context.Context,
	projectID string,
	sprintID *string,
	assigneeID *string,
	labelIDs []string,
) *apperrors.AppError {
	projectExists, err := s.projects.ExistsByID(ctx, projectID)
	if err != nil {
		return apperrors.Internal("failed to validate project")
	}
	if !projectExists {
		return apperrors.NotFound("project not found")
	}

	if sprintID != nil {
		sprint, err := s.sprints.FindByID(ctx, *sprintID)
		if err != nil {
			if errors.Is(err, repositories.ErrNotFound) {
				return apperrors.NotFound("sprint not found")
			}
			return apperrors.Internal("failed to validate sprint")
		}
		if sprint.ProjectID != projectID {
			return apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{
				"sprint_id": "Select a sprint that belongs to the selected project.",
			})
		}
	}

	if assigneeID != nil {
		exists, err := s.users.ExistsByID(ctx, *assigneeID)
		if err != nil {
			return apperrors.Internal("failed to validate assignee")
		}
		if !exists {
			return apperrors.NotFound("assignee not found")
		}
	}

	if len(labelIDs) > 0 {
		distinct := make(map[string]struct{}, len(labelIDs))
		for _, id := range labelIDs {
			if _, ok := distinct[id]; ok {
				return apperrors.Validation("Please correct the highlighted fields and try again.", apperrors.FieldErrors{
					"label_ids": "Remove duplicate labels.",
				})
			}
			distinct[id] = struct{}{}
		}
		exists, err := s.labels.ExistsByIDs(ctx, labelIDs)
		if err != nil {
			return apperrors.Internal("failed to validate labels")
		}
		if !exists {
			return apperrors.NotFound("one or more labels not found")
		}
	}

	return nil
}

func (s *IssueService) hydrateIssues(ctx context.Context, rows []models.Issue) ([]IssueSummary, *apperrors.AppError) {
	if len(rows) == 0 {
		return []IssueSummary{}, nil
	}

	projectIDs := make([]string, 0, len(rows))
	sprintIDs := make([]string, 0, len(rows))
	userIDs := make([]string, 0, len(rows)*2)
	issueIDs := make([]string, 0, len(rows))
	projectSet := map[string]struct{}{}
	sprintSet := map[string]struct{}{}
	userSet := map[string]struct{}{}

	for _, issue := range rows {
		issueIDs = append(issueIDs, issue.ID)
		if _, ok := projectSet[issue.ProjectID]; !ok {
			projectSet[issue.ProjectID] = struct{}{}
			projectIDs = append(projectIDs, issue.ProjectID)
		}
		if issue.SprintID != nil {
			if _, ok := sprintSet[*issue.SprintID]; !ok {
				sprintSet[*issue.SprintID] = struct{}{}
				sprintIDs = append(sprintIDs, *issue.SprintID)
			}
		}
		if _, ok := userSet[issue.CreatedBy]; !ok {
			userSet[issue.CreatedBy] = struct{}{}
			userIDs = append(userIDs, issue.CreatedBy)
		}
		if issue.AssigneeID != nil {
			if _, ok := userSet[*issue.AssigneeID]; !ok {
				userSet[*issue.AssigneeID] = struct{}{}
				userIDs = append(userIDs, *issue.AssigneeID)
			}
		}
	}

	projectsMap, err := s.projects.SummariesByIDs(ctx, projectIDs)
	if err != nil {
		return nil, apperrors.Internal("failed to load issue projects")
	}
	sprintsMap, err := s.sprints.SummariesByIDs(ctx, sprintIDs)
	if err != nil {
		return nil, apperrors.Internal("failed to load issue sprints")
	}
	usersRows, err := s.users.FindByIDs(ctx, userIDs)
	if err != nil {
		return nil, apperrors.Internal("failed to load issue users")
	}
	userMap := make(map[string]UserSummary, len(usersRows))
	for _, user := range usersRows {
		userMap[user.ID] = UserSummary{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	labelsMapRaw, err := s.issues.ListLabelsByIssueIDs(ctx, issueIDs)
	if err != nil {
		return nil, apperrors.Internal("failed to load issue labels")
	}
	labelsMap := make(map[string][]LabelSummary, len(labelsMapRaw))
	for issueID, labels := range labelsMapRaw {
		items := make([]LabelSummary, 0, len(labels))
		for _, label := range labels {
			items = append(items, LabelSummary{
				ID:          label.ID,
				Name:        label.Name,
				Color:       label.Color,
				Description: label.Description,
				CreatedAt:   label.CreatedAt,
				UpdatedAt:   label.UpdatedAt,
			})
		}
		labelsMap[issueID] = items
	}

	items := make([]IssueSummary, 0, len(rows))
	for _, issue := range rows {
		projectRow := projectsMap[issue.ProjectID]
		project := ProjectSummary{
			ID:          projectRow.ID,
			Name:        projectRow.Name,
			Description: projectRow.Description,
			Key:         projectRow.Key,
			CreatedBy:   projectRow.CreatedBy,
			CreatedAt:   projectRow.CreatedAt,
			UpdatedAt:   projectRow.UpdatedAt,
			IssueCounts: mapIssueCounts(projectRow.IssueCounts),
		}
		if projectRow.ActiveSprint != nil {
			project.ActiveSprint = mapSprintRow(*projectRow.ActiveSprint)
		}

		var sprint *SprintSummary
		if issue.SprintID != nil {
			if sprintRow, ok := sprintsMap[*issue.SprintID]; ok {
				sprint = mapSprintRow(sprintRow)
			}
		}

		var assignee *UserSummary
		if issue.AssigneeID != nil {
			if u, ok := userMap[*issue.AssigneeID]; ok {
				user := u
				assignee = &user
			}
		}

		creator := userMap[issue.CreatedBy]
		labels := labelsMap[issue.ID]
		if labels == nil {
			labels = []LabelSummary{}
		}

		items = append(items, IssueSummary{
			ID:          issue.ID,
			Identifier:  issue.Identifier,
			Title:       issue.Title,
			Description: issue.Description,
			Status:      issue.Status,
			Priority:    issue.Priority,
			ProjectID:   issue.ProjectID,
			SprintID:    issue.SprintID,
			AssigneeID:  issue.AssigneeID,
			CreatedBy:   issue.CreatedBy,
			ArchivedAt:  issue.ArchivedAt,
			ArchivedBy:  issue.ArchivedBy,
			CreatedAt:   issue.CreatedAt,
			UpdatedAt:   issue.UpdatedAt,
			Project:     project,
			Sprint:      sprint,
			Assignee:    assignee,
			Creator:     creator,
			Labels:      labels,
		})
	}

	return items, nil
}

func (s *IssueService) hydrateIssueActivities(ctx context.Context, issueID string) ([]IssueActivity, error) {
	rows, err := s.issues.ListActivitiesByIssueID(ctx, issueID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []IssueActivity{}, nil
	}

	userIDs := make([]string, 0, len(rows))
	seen := map[string]struct{}{}
	for _, row := range rows {
		if _, ok := seen[row.UserID]; !ok {
			seen[row.UserID] = struct{}{}
			userIDs = append(userIDs, row.UserID)
		}
	}
	users, err := s.users.FindByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]UserSummary, len(users))
	for _, user := range users {
		userMap[user.ID] = UserSummary{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	out := make([]IssueActivity, 0, len(rows))
	for _, row := range rows {
		var user *UserSummary
		if mapped, ok := userMap[row.UserID]; ok {
			copy := mapped
			user = &copy
		}
		out = append(out, IssueActivity{
			ID:        row.ID,
			IssueID:   row.IssueID,
			UserID:    row.UserID,
			Action:    row.Action,
			FieldName: row.FieldName,
			OldValue:  row.OldValue,
			NewValue:  row.NewValue,
			CreatedAt: row.CreatedAt,
			User:      user,
		})
	}

	return out, nil
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}
