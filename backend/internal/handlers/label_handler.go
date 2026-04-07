package handlers

import (
	"context"
	"net/http"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/abhinavmaity/linear-lite/backend/internal/validation"
	"github.com/gin-gonic/gin"
)

type LabelService interface {
	List(ctx context.Context, input services.LabelListInput) ([]services.LabelSummary, int64, *apperrors.AppError)
}

type LabelHandler struct {
	service LabelService
}

func NewLabelHandler(service LabelService) *LabelHandler {
	return &LabelHandler{service: service}
}

func (h *LabelHandler) List(c *gin.Context) {
	pagination, appErr := validation.ParsePagination(c, 100, 100)
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	sortBy, appErr := validation.ParseSortField(c.Query("sort_by"), "name", []string{"name", "created_at"})
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}
	sortOrder, appErr := validation.ParseSortOrder(c.Query("sort_order"))
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	items, total, serviceErr := h.service.List(c, services.LabelListInput{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
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
