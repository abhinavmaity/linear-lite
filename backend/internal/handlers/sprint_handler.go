package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/abhinavmaity/linear-lite/backend/internal/validation"
	"github.com/gin-gonic/gin"
)

type SprintService interface {
	List(ctx context.Context, input services.SprintListInput) ([]services.SprintSummary, int64, *apperrors.AppError)
	Create(ctx context.Context, input services.SprintCreateInput) (*services.SprintDetail, *apperrors.AppError)
	Get(ctx context.Context, id string) (*services.SprintDetail, *apperrors.AppError)
	Update(ctx context.Context, id string, input services.SprintUpdateInput) (*services.SprintDetail, *apperrors.AppError)
	Delete(ctx context.Context, id string) *apperrors.AppError
}

type SprintHandler struct {
	service SprintService
}

func NewSprintHandler(service SprintService) *SprintHandler {
	return &SprintHandler{service: service}
}

func (h *SprintHandler) List(c *gin.Context) {
	pagination, appErr := validation.ParsePagination(c, 50, 100)
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	sortBy, appErr := validation.ParseSortField(c.Query("sort_by"), "start_date", []string{"name", "start_date", "end_date", "created_at"})
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}
	sortOrder, appErr := validation.ParseSortOrder(c.Query("sort_order"))
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	projectID, appErr := validation.ParseOptionalUUIDQuery(c, "project_id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	rawStatus := c.Query("status")
	if enumErr := validation.ValidateEnum("status", rawStatus, []string{"planned", "active", "completed"}); enumErr != nil {
		apperrors.Write(c, enumErr, requestID(c))
		return
	}
	var status *string
	if rawStatus != "" {
		status = &rawStatus
	}

	var projectIDString *string
	if projectID != nil {
		id := projectID.String()
		projectIDString = &id
	}

	items, total, serviceErr := h.service.List(c, services.SprintListInput{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		ProjectID: projectIDString,
		Status:    status,
		Search:    c.Query("search"),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteCollection(c, http.StatusOK, items, BuildPaginationMeta(pagination.Page, pagination.Limit, total))
}

type createSprintRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ProjectID   string  `json:"project_id"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Status      *string `json:"status"`
}

type updateSprintRequest struct {
	Name        *string  `json:"name"`
	Description **string `json:"description"`
	StartDate   *string  `json:"start_date"`
	EndDate     *string  `json:"end_date"`
	Status      *string  `json:"status"`
}

func (h *SprintHandler) Create(c *gin.Context) {
	var req createSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			"body": "must be valid JSON",
		}), requestID(c))
		return
	}

	sprint, serviceErr := h.service.Create(c, services.SprintCreateInput{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      req.Status,
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusCreated, sprint)
}

func (h *SprintHandler) Get(c *gin.Context) {
	id, appErr := validation.ParseUUIDParam(c, "id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	sprint, serviceErr := h.service.Get(c, id.String())
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, sprint)
}

func (h *SprintHandler) Update(c *gin.Context) {
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

	var req updateSprintRequest
	if v, ok, err := parseStringPointer(raw, "name"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.Name = v
	}
	if v, ok, err := parseNullableStringPointer(raw, "description"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.Description = &v
	}
	if v, ok, err := parseStringPointer(raw, "start_date"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.StartDate = v
	}
	if v, ok, err := parseStringPointer(raw, "end_date"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.EndDate = v
	}
	if v, ok, err := parseStringPointer(raw, "status"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.Status = v
	}

	sprint, serviceErr := h.service.Update(c, id.String(), services.SprintUpdateInput{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      req.Status,
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, sprint)
}

func (h *SprintHandler) Delete(c *gin.Context) {
	id, appErr := validation.ParseUUIDParam(c, "id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	serviceErr := h.service.Delete(c, id.String())
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	c.Status(http.StatusNoContent)
}
