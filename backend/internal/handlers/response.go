package handlers

import (
	"github.com/abhinavmaity/linear-lite/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type CollectionResponse struct {
	Items      any `json:"items"`
	Pagination any `json:"pagination"`
}

type ResourceResponse struct {
	Data any `json:"data"`
}

func WriteCollection(c *gin.Context, status int, items any, pagination any) {
	c.JSON(status, CollectionResponse{
		Items:      items,
		Pagination: pagination,
	})
}

func WriteResource(c *gin.Context, status int, data any) {
	c.JSON(status, ResourceResponse{
		Data: data,
	})
}

func requestID(c *gin.Context) string {
	return middleware.GetRequestID(c)
}
