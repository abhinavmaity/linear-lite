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

type LabelService interface {
	List(ctx context.Context, input services.LabelListInput) ([]services.LabelSummary, int64, *apperrors.AppError)
	Create(ctx context.Context, input services.LabelCreateInput) (*services.LabelSummary, *apperrors.AppError)
	Get(ctx context.Context, id string) (*services.LabelDetail, *apperrors.AppError)
	Update(ctx context.Context, id string, input services.LabelUpdateInput) (*services.LabelSummary, *apperrors.AppError)
	Delete(ctx context.Context, id string) *apperrors.AppError
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

type createLabelRequest struct {
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Description *string `json:"description"`
}

type updateLabelRequest struct {
	Name        *string  `json:"name"`
	Color       *string  `json:"color"`
	Description **string `json:"description"`
}

func (h *LabelHandler) Create(c *gin.Context) {
	var req createLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			"body": "must be valid JSON",
		}), requestID(c))
		return
	}

	label, serviceErr := h.service.Create(c, services.LabelCreateInput{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusCreated, label)
}

func (h *LabelHandler) Get(c *gin.Context) {
	id, appErr := validation.ParseUUIDParam(c, "id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	label, serviceErr := h.service.Get(c, id.String())
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, label)
}

func (h *LabelHandler) Update(c *gin.Context) {
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

	var req updateLabelRequest
	if v, ok, err := parseStringPointer(raw, "name"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.Name = v
	}
	if v, ok, err := parseStringPointer(raw, "color"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.Color = v
	}
	if v, ok, err := parseNullableStringPointer(raw, "description"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.Description = &v
	}

	label, serviceErr := h.service.Update(c, id.String(), services.LabelUpdateInput{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, label)
}

func (h *LabelHandler) Delete(c *gin.Context) {
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
