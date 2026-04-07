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
	projectRepo := repositories.NewProjectRepository(deps.DB)
	sprintRepo := repositories.NewSprintRepository(deps.DB)
	labelRepo := repositories.NewLabelRepository(deps.DB)
	issueRepo := repositories.NewIssueRepository(deps.DB)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTTTL, cfg.BcryptCost)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo)
	sprintService := services.NewSprintService(sprintRepo)
	labelService := services.NewLabelService(labelRepo)
	issueService := services.NewIssueService(issueRepo, userRepo, projectRepo, sprintRepo, labelRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	sprintHandler := handlers.NewSprintHandler(sprintService)
	labelHandler := handlers.NewLabelHandler(labelService)
	issueHandler := handlers.NewIssueHandler(issueService)

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.CORSOrigins))

	registerRoutes(router, cfg, authHandler, userHandler, projectHandler, sprintHandler, labelHandler, issueHandler)

	return router
}

func registerRoutes(
	router *gin.Engine,
	cfg config.Config,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	projectHandler *handlers.ProjectHandler,
	sprintHandler *handlers.SprintHandler,
	labelHandler *handlers.LabelHandler,
	issueHandler *handlers.IssueHandler,
) {
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
	protected.GET("/users", userHandler.List)
	protected.GET("/projects", projectHandler.List)
	protected.GET("/sprints", sprintHandler.List)
	protected.GET("/labels", labelHandler.List)
	protected.GET("/issues", issueHandler.List)
	protected.POST("/issues", issueHandler.Create)
	protected.GET("/issues/:id", issueHandler.Get)
	protected.PUT("/issues/:id", issueHandler.Update)
	protected.DELETE("/issues/:id", issueHandler.Delete)
}
