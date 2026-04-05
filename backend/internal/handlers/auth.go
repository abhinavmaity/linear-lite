package handlers

import (
	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	apperrors.Write(c, apperrors.NotImplemented("register handler is not implemented yet"), requestID(c))
}

func Login(c *gin.Context) {
	apperrors.Write(c, apperrors.NotImplemented("login handler is not implemented yet"), requestID(c))
}

func Me(c *gin.Context) {
	apperrors.Write(c, apperrors.NotImplemented("me handler is not implemented yet"), requestID(c))
}
