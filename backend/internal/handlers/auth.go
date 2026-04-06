package handlers

import (
	"context"
	"net/http"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/abhinavmaity/linear-lite/backend/internal/middleware"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx context.Context, input services.RegisterInput) (*services.AuthSession, *apperrors.AppError)
	Login(ctx context.Context, input services.LoginInput) (*services.AuthSession, *apperrors.AppError)
	Me(ctx context.Context, userID string) (*services.AuthUser, *apperrors.AppError)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			"body": "must be valid JSON",
		}), requestID(c))
		return
	}

	session, appErr := h.service.Register(c, services.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusCreated, session)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.Validation("invalid request body", apperrors.FieldErrors{
			"body": "must be valid JSON",
		}), requestID(c))
		return
	}

	session, appErr := h.service.Login(c, services.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, session)
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, appErr := h.service.Me(c, middleware.AuthUserID(c))
	if appErr != nil {
		apperrors.Write(c, appErr, requestID(c))
		return
	}

	WriteResource(c, http.StatusOK, user)
}
