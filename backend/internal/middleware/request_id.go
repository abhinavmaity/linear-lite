package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextKeyRequestID = "request_id"
	HeaderRequestID     = "X-Request-ID"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(HeaderRequestID))
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(ContextKeyRequestID, requestID)
		c.Writer.Header().Set(HeaderRequestID, requestID)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	raw, exists := c.Get(ContextKeyRequestID)
	if !exists {
		return ""
	}

	requestID, ok := raw.(string)
	if !ok {
		return ""
	}

	return requestID
}
