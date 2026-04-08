package handlers

import (
	"context"
	"net/http"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/middleware"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type DashboardService interface {
	GetStats(ctx context.Context, userID string) (*services.DashboardStats, *apperrors.AppError)
}

type DashboardHandler struct {
	service DashboardService
}

func NewDashboardHandler(service DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	stats, appErr := h.service.GetStats(c, middleware.AuthUserID(c))
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, stats)
}
