package server

import (
	"net/http"

	"github.com/abhinavmaity/linear-lite/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func New(cfg config.Config, deps Dependencies) *gin.Engine {
	_ = cfg
	_ = deps

	router := gin.New()
	router.Use(gin.Recovery())

	registerRoutes(router)

	return router
}

func registerRoutes(router *gin.Engine) {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
