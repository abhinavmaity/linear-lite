package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/middleware"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/abhinavmaity/linear-lite/backend/internal/validation"
	"github.com/gin-gonic/gin"
)

type ProjectService interface {
	List(ctx context.Context, input services.ProjectListInput) ([]services.ProjectSummary, int64, *apperrors.AppError)
	Create(ctx context.Context, actorID string, input services.ProjectCreateInput) (*services.ProjectDetail, *apperrors.AppError)
	Get(ctx context.Context, id string) (*services.ProjectDetail, *apperrors.AppError)
	Update(ctx context.Context, id string, input services.ProjectUpdateInput) (*services.ProjectDetail, *apperrors.AppError)
	Delete(ctx context.Context, id string) *apperrors.AppError
}

type ProjectHandler struct {
	service ProjectService
}

func NewProjectHandler(service ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) List(c *gin.Context) {
	pagination, appErr := validation.ParsePagination(c, 50, 100)
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	sortBy, appErr := validation.ParseSortField(c.Query("sort_by"), "name", []string{"name", "created_at", "updated_at"})
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}
	sortOrder, appErr := validation.ParseSortOrder(c.Query("sort_order"))
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	items, total, serviceErr := h.service.List(c, services.ProjectListInput{
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

type createProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Key         string  `json:"key"`
}

type updateProjectRequest struct {
	Name        *string  `json:"name"`
	Description **string `json:"description"`
	Key         *string  `json:"key"`
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.Validation("Request body is invalid.", apperrors.FieldErrors{
			"body": "Body must be valid JSON.",
		}), requestID(c))
		return
	}

	project, serviceErr := h.service.Create(c, middleware.AuthUserID(c), services.ProjectCreateInput{
		Name:        req.Name,
		Description: req.Description,
		Key:         req.Key,
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusCreated, project)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, appErr := validation.ParseUUIDParam(c, "id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	project, serviceErr := h.service.Get(c, id.String())
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, appErr := validation.ParseUUIDParam(c, "id")
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	var raw map[string]json.RawMessage
	if err := c.ShouldBindJSON(&raw); err != nil {
		apperrors.Write(c, apperrors.Validation("Request body is invalid.", apperrors.FieldErrors{
			"body": "Body must be valid JSON.",
		}), requestID(c))
		return
	}

	var req updateProjectRequest
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
	if v, ok, err := parseStringPointer(raw, "key"); err != nil {
		apperrors.Write(c, err, requestID(c))
		return
	} else if ok {
		req.Key = v
	}

	project, serviceErr := h.service.Update(c, id.String(), services.ProjectUpdateInput{
		Name:        req.Name,
		Description: req.Description,
		Key:         req.Key,
	})
	if serviceErr != nil {
		apperrors.Write(c, serviceErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
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
