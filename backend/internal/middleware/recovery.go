package middleware

import (
	"log/slog"
	"runtime/debug"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error("panic recovered",
					"request_id", GetRequestID(c),
					"panic", recovered,
					"stack", string(debug.Stack()),
				)

				apperrors.Write(c, apperrors.Internal("unexpected server error"), GetRequestID(c))
			}
		}()

		c.Next()
	}
}
