package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/middleware"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/abhinavmaity/linear-lite/backend/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IssueService interface {
	List(ctx context.Context, input services.IssueListInput) ([]services.IssueSummary, int64, *apperrors.AppError)
	Get(ctx context.Context, id string, includeArchived bool) (*services.IssueDetail, *apperrors.AppError)
	Create(ctx context.Context, actorID string, input services.CreateIssueInput) (*services.IssueDetail, *apperrors.AppError)
	Update(ctx context.Context, actorID string, input services.UpdateIssueInput) (*services.IssueDetail, *apperrors.AppError)
	Archive(ctx context.Context, actorID string, id string) *apperrors.AppError
}

type IssueHandler struct {
	service IssueService
}

func NewIssueHandler(service IssueService) *IssueHandler {
	return &IssueHandler{service: service}
}

type createIssueRequest struct {
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	Status      *string  `json:"status"`
	Priority    *string  `json:"priority"`
	ProjectID   string   `json:"project_id"`
	SprintID    *string  `json:"sprint_id"`
	AssigneeID  *string  `json:"assignee_id"`
	LabelIDs    []string `json:"label_ids"`
}

func (h *IssueHandler) List(c *gin.Context) {
	pagination, appErr := validation.ParsePagination(c, 50, 100)
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	sortBy, appErr := validation.ParseSortField(c.Query("sort_by"), "updated_at", []string{"identifier", "title", "status", "priority", "created_at", "updated_at"})
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}
	sortOrder, appErr := validation.ParseSortOrder(c.Query("sort_order"))
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	statuses := validation.ParseRepeatedQuery(c, "status")
	for _, status := range statuses {
		if err := validation.ValidateEnum("status", status, []string{"backlog", "todo", "in_progress", "in_review", "done", "cancelled"}); err != nil {
			apperrors.Write(c, err, requestID(c))
			return
		}
	}
	priorities := validation.ParseRepeatedQuery(c, "priority")
	for _, priority := range priorities {
		if err := validation.ValidateEnum("priority", priority, []string{"low", "medium", "high", "urgent"}); err != nil {
			apperrors.Write(c, err, requestID(c))
			return
		}
	}

	assigneeID, appErr := validation.ParseOptionalUUIDQuery(c, "assignee_id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}
	projectID, appErr := validation.ParseOptionalUUIDQuery(c, "project_id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}
	sprintID, appErr := validation.ParseOptionalUUIDQuery(c, "sprint_id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	labelIDRaw := validation.ParseRepeatedQuery(c, "label_id")
	labelIDsParsed, appErr := validation.ParseDistinctUUIDArray("label_id", labelIDRaw)
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}
	labelMode := strings.TrimSpace(c.DefaultQuery("label_mode", "any"))
	if labelMode != "any" && labelMode != "all" {
		apperrors.Write(c, apperrors.Validation("invalid query parameter", apperrors.FieldErrors{
			"label_mode": "must be any or all",
		}), requestID(c))
		return
	}
	includeArchived, appErr := validation.ParseOptionalBoolQuery(c, "include_archived", false)
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	items, total, serviceErr := h.service.List(c, services.IssueListInput{
		Page:            pagination.Page,
		Limit:           pagination.Limit,
		SortBy:          sortBy,
		SortOrder:       sortOrder,
		Search:          c.Query("search"),
		Statuses:        statuses,
		Priorities:      priorities,
		AssigneeID:      uuidPtrToString(assigneeID),
		ProjectID:       uuidPtrToString(projectID),
		SprintID:        uuidPtrToString(sprintID),
		LabelIDs:        uuidSliceToString(labelIDsParsed),
		LabelMode:       labelMode,
		IncludeArchived: includeArchived,
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteCollection(c, http.StatusOK, items, BuildPaginationMeta(pagination.Page, pagination.Limit, total))
}

func (h *IssueHandler) Create(c *gin.Context) {
	var req createIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			"body": "must be valid JSON",
		}), requestID(c))
		return
	}

	if _, err := uuid.Parse(req.ProjectID); err != nil {
		apperrors.Write(c, apperrors.Validation("one or more fields are invalid", apperrors.FieldErrors{
			"project_id": "must be a valid UUID",
		}), requestID(c))
		return
	}
	if req.SprintID != nil {
		if _, err := uuid.Parse(*req.SprintID); err != nil {
			apperrors.Write(c, apperrors.Validation("one or more fields are invalid", apperrors.FieldErrors{
				"sprint_id": "must be a valid UUID",
			}), requestID(c))
			return
		}
	}
	if req.AssigneeID != nil {
		if _, err := uuid.Parse(*req.AssigneeID); err != nil {
			apperrors.Write(c, apperrors.Validation("one or more fields are invalid", apperrors.FieldErrors{
				"assignee_id": "must be a valid UUID",
			}), requestID(c))
			return
		}
	}
	if req.Status != nil {
		if err := validation.ValidateEnum("status", *req.Status, []string{"backlog", "todo", "in_progress", "in_review", "done", "cancelled"}); err != nil {
			apperrors.Write(c, err, requestID(c))
			return
		}
	}
	if req.Priority != nil {
		if err := validation.ValidateEnum("priority", *req.Priority, []string{"low", "medium", "high", "urgent"}); err != nil {
			apperrors.Write(c, err, requestID(c))
			return
		}
	}

	labelUUIDs, appErr := validation.ParseDistinctUUIDArray("label_ids", req.LabelIDs)
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	issue, serviceErr := h.service.Create(c, middleware.AuthUserID(c), services.CreateIssueInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		ProjectID:   req.ProjectID,
		SprintID:    req.SprintID,
		AssigneeID:  req.AssigneeID,
		LabelIDs:    uuidSliceToString(labelUUIDs),
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusCreated, issue)
}

func (h *IssueHandler) Get(c *gin.Context) {
	id, appErr := validation.ParseUUIDParam(c, "id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	includeArchived, appErr := validation.ParseOptionalBoolQuery(c, "include_archived", false)
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	issue, serviceErr := h.service.Get(c, id.String(), includeArchived)
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, issue)
}

func (h *IssueHandler) Update(c *gin.Context) {
	id, appErr := validation.ParseUUIDParam(c, "id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	var raw map[string]json.RawMessage
	if err := c.ShouldBindJSON(&raw); err != nil {
		apperrors.Write(c, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			"body": "must be valid JSON",
		}), requestID(c))
		return
	}

	var input services.UpdateIssueInput
	input.ID = id.String()

	if v, ok, err := parseRequiredStringPointer(raw, "title"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		input.Title = v
	}
	if v, ok, err := parseNullableStringPointer(raw, "description"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		input.Description = &v
	}
	if v, ok, err := parseStringPointer(raw, "status"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		if enumErr := validation.ValidateEnum("status", *v, []string{"backlog", "todo", "in_progress", "in_review", "done", "cancelled"}); enumErr != nil {
			apperrors.Write(c, enumErr, requestID(c))
			return
		}
		input.Status = v
	}
	if v, ok, err := parseStringPointer(raw, "priority"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		if enumErr := validation.ValidateEnum("priority", *v, []string{"low", "medium", "high", "urgent"}); enumErr != nil {
			apperrors.Write(c, enumErr, requestID(c))
			return
		}
		input.Priority = v
	}
	if v, ok, err := parseUUIDStringPointer(raw, "project_id", false); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		input.ProjectID = v
	}
	if v, ok, err := parseUUIDStringPointer(raw, "sprint_id", true); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		input.SprintID = &v
	}
	if v, ok, err := parseUUIDStringPointer(raw, "assignee_id", true); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		input.AssigneeID = &v
	}
	if v, ok, err := parseUUIDStringSlice(raw, "label_ids"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		input.LabelIDs = &v
	}
	if v, ok, err := parseBoolPointer(raw, "archived"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		input.Archived = v
	}

	issue, serviceErr := h.service.Update(c, middleware.AuthUserID(c), input)
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, issue)
}

func (h *IssueHandler) Delete(c *gin.Context) {
	id, appErr := validation.ParseUUIDParam(c, "id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	if err := h.service.Archive(c, middleware.AuthUserID(c), id.String()); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	}

	c.Status(http.StatusNoContent)
}

func uuidPtrToString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	value := id.String()
	return &value
}

func uuidSliceToString(ids []uuid.UUID) []string {
	if len(ids) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, id.String())
	}
	return out
}

func parseStringPointer(raw map[string]json.RawMessage, key string) (*string, bool, *apperrors.AppError) {
	v, ok := raw[key]
	if !ok {
		return nil, false, nil
	}
	var value string
	if err := json.Unmarshal(v, &value); err != nil {
		return nil, false, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			key: "must be a string",
		})
	}
	value = strings.TrimSpace(value)
	return &value, true, nil
}

func parseRequiredStringPointer(raw map[string]json.RawMessage, key string) (*string, bool, *apperrors.AppError) {
	value, ok, err := parseStringPointer(raw, key)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	if value == nil || *value == "" {
		return nil, false, apperrors.Validation("one or more fields are invalid", apperrors.FieldErrors{
			key: "is required",
		})
	}
	return value, true, nil
}

func parseNullableStringPointer(raw map[string]json.RawMessage, key string) (*string, bool, *apperrors.AppError) {
	v, ok := raw[key]
	if !ok {
		return nil, false, nil
	}
	if string(v) == "null" {
		return nil, true, nil
	}
	var value string
	if err := json.Unmarshal(v, &value); err != nil {
		return nil, false, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			key: "must be a string or null",
		})
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, true, nil
	}
	return &value, true, nil
}

func parseUUIDStringPointer(raw map[string]json.RawMessage, key string, nullable bool) (*string, bool, *apperrors.AppError) {
	v, ok := raw[key]
	if !ok {
		return nil, false, nil
	}
	if nullable && string(v) == "null" {
		return nil, true, nil
	}
	var value string
	if err := json.Unmarshal(v, &value); err != nil {
		msg := "must be a UUID string"
		if nullable {
			msg = "must be a UUID string or null"
		}
		return nil, false, apperrors.Validation("invalid request body", apperrors.FieldErrors{key: msg})
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, false, apperrors.Validation("one or more fields are invalid", apperrors.FieldErrors{
			key: "must be a valid UUID",
		})
	}
	if _, err := uuid.Parse(value); err != nil {
		return nil, false, apperrors.Validation("one or more fields are invalid", apperrors.FieldErrors{
			key: "must be a valid UUID",
		})
	}
	return &value, true, nil
}

func parseUUIDStringSlice(raw map[string]json.RawMessage, key string) ([]string, bool, *apperrors.AppError) {
	v, ok := raw[key]
	if !ok {
		return nil, false, nil
	}
	var value []string
	if err := json.Unmarshal(v, &value); err != nil {
		return nil, false, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			key: "must be an array of UUID strings",
		})
	}
	parsed, err := validation.ParseDistinctUUIDArray(key, value)
	if err != nil {
		return nil, false, err
	}
	return uuidSliceToString(parsed), true, nil
}

func parseBoolPointer(raw map[string]json.RawMessage, key string) (*bool, bool, *apperrors.AppError) {
	v, ok := raw[key]
	if !ok {
		return nil, false, nil
	}
	var value bool
	if err := json.Unmarshal(v, &value); err != nil {
		return nil, false, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			key: "must be a boolean",
		})
	}
	return &value, true, nil
}
