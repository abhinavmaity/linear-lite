package server

import (
	"net/http"

	cachepkg "github.com/abhinavmaity/linear-lite/backend/internal/cache"
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
	cacheStore := cachepkg.NewStore(deps.Redis)

	userRepo := repositories.NewUserRepository(deps.DB)
	projectRepo := repositories.NewProjectRepository(deps.DB)
	sprintRepo := repositories.NewSprintRepository(deps.DB)
	labelRepo := repositories.NewLabelRepository(deps.DB)
	issueRepo := repositories.NewIssueRepository(deps.DB)
	googleVerifier := services.NewGoogleIDTokenVerifier()

	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTTTL, cfg.BcryptCost, cfg.GoogleClientID, googleVerifier, cacheStore)
	userService := services.NewUserService(userRepo, cacheStore)
	projectService := services.NewProjectService(projectRepo, userRepo, cacheStore)
	sprintService := services.NewSprintService(sprintRepo, projectRepo, cacheStore)
	labelService := services.NewLabelService(labelRepo, cacheStore)
	issueService := services.NewIssueService(issueRepo, userRepo, projectRepo, sprintRepo, labelRepo, cacheStore)
	dashboardService := services.NewDashboardService(issueRepo, sprintRepo, userRepo, cacheStore)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	sprintHandler := handlers.NewSprintHandler(sprintService)
	labelHandler := handlers.NewLabelHandler(labelService)
	issueHandler := handlers.NewIssueHandler(issueService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.CORSOrigins))

	registerRoutes(router, cfg, authHandler, userHandler, projectHandler, sprintHandler, labelHandler, issueHandler, dashboardHandler)

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
	dashboardHandler *handlers.DashboardHandler,
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
	public.POST("/auth/google", authHandler.LoginWithGoogle)

	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(cfg.JWTSecret))
	protected.GET("/auth/me", authHandler.Me)
	protected.GET("/users", userHandler.List)
	protected.GET("/users/:id", userHandler.Get)
	protected.GET("/projects", projectHandler.List)
	protected.POST("/projects", projectHandler.Create)
	protected.GET("/projects/:id", projectHandler.Get)
	protected.PUT("/projects/:id", projectHandler.Update)
	protected.DELETE("/projects/:id", projectHandler.Delete)
	protected.GET("/sprints", sprintHandler.List)
	protected.POST("/sprints", sprintHandler.Create)
	protected.GET("/sprints/:id", sprintHandler.Get)
	protected.PUT("/sprints/:id", sprintHandler.Update)
	protected.DELETE("/sprints/:id", sprintHandler.Delete)
	protected.GET("/labels", labelHandler.List)
	protected.POST("/labels", labelHandler.Create)
	protected.GET("/labels/:id", labelHandler.Get)
	protected.PUT("/labels/:id", labelHandler.Update)
	protected.DELETE("/labels/:id", labelHandler.Delete)
	protected.GET("/dashboard/stats", dashboardHandler.Stats)
	protected.GET("/issues", issueHandler.List)
	protected.POST("/issues", issueHandler.Create)
	protected.GET("/issues/:id", issueHandler.Get)
	protected.PUT("/issues/:id", issueHandler.Update)
	protected.DELETE("/issues/:id", issueHandler.Delete)
}
