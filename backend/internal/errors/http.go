package errors

import "github.com/gin-gonic/gin"

func Write(c *gin.Context, err *AppError, requestID string) {
	c.AbortWithStatusJSON(err.Status, err.Response(requestID))
}
