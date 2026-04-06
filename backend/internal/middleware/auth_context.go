package middleware

import "github.com/gin-gonic/gin"

func AuthUserID(c *gin.Context) string {
	raw, exists := c.Get(ContextKeyAuthUserID)
	if !exists {
		return ""
	}
	value, ok := raw.(string)
	if !ok {
		return ""
	}
	return value
}

func AuthEmail(c *gin.Context) string {
	raw, exists := c.Get(ContextKeyAuthEmail)
	if !exists {
		return ""
	}
	value, ok := raw.(string)
	if !ok {
		return ""
	}
	return value
}
