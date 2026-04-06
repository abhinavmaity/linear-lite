package server

import (
	"net/http"

	"github.com/abhinavmaity/linear-lite/backend/internal/config"
	"github.com/abhinavmaity/linear-lite/backend/internal/handlers"
	"github.com/abhinavmaity/linear-lite/backend/internal/middleware"
	"github.com/abhinavmaity/linear-lite/backend/internal/repositories"
	"github.com/abhinavmaity/linear-lite/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func New(cfg config.Config, deps Dependencies) *gin.Engine {
	userRepo := repositories.NewUserRepository(deps.DB)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTTTL, cfg.BcryptCost)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.CORSOrigins))

	registerRoutes(router, cfg, authHandler)

	return router
}

func registerRoutes(router *gin.Engine, cfg config.Config, authHandler *handlers.AuthHandler) {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	v1 := router.Group("/api/v1")

	public := v1.Group("")
	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)

	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(cfg.JWTSecret))
	protected.GET("/auth/me", authHandler.Me)
}
