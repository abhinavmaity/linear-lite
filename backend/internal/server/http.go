package server

import (
	"net/http"

	"github.com/abhinavmaity/linear-lite/backend/internal/config"
	"github.com/abhinavmaity/linear-lite/backend/internal/handlers"
	"github.com/abhinavmaity/linear-lite/backend/internal/middleware"
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

	registerRoutes(router, cfg)

	return router
}

func registerRoutes(router *gin.Engine, cfg config.Config) {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	v1 := router.Group("/api/v1")

	public := v1.Group("")
	public.POST("/auth/register", handlers.Register)
	public.POST("/auth/login", handlers.Login)

	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(cfg.JWTSecret))
	protected.GET("/auth/me", handlers.Me)
}
