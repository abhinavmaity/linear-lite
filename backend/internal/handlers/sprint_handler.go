package handlers

import (
	"context"
	"net/http"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/abhinavmaity/linear-lite/backend/internal/validation"
	"github.com/gin-gonic/gin"
)

type SprintService interface {
	List(ctx context.Context, input services.SprintListInput) ([]services.SprintSummary, int64, *apperrors.AppError)
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
