package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/abhinavmaity/linear-lite/backend/internal/cache"
	"github.com/abhinavmaity/linear-lite/backend/internal/config"
	"github.com/abhinavmaity/linear-lite/backend/internal/database"
	"github.com/abhinavmaity/linear-lite/backend/internal/server"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	cfg        config.Config
	httpServer *http.Server
	db         *gorm.DB
	redis      *redis.Client
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize postgres: %w", err)
	}

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize redis: %w", err)
	}

	router := server.New(cfg, server.Dependencies{
		DB:    db,
		Redis: redisClient,
	})

	httpServer := &http.Server{
		Addr:              cfg.ListenAddr(),
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
		ReadTimeout:       cfg.HTTPReadTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}

	return &App{
		cfg:        cfg,
		httpServer: httpServer,
		db:         db,
		redis:      redisClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		slog.Info("starting api server", "env", a.cfg.AppEnv, "addr", a.cfg.ListenAddr())
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		return a.shutdown()
	case err := <-errCh:
		_ = a.shutdown()
		return fmt.Errorf("server failure: %w", err)
	}
}

func (a *App) shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTPShutdownTimeout)
	defer cancel()

	var joined error

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		joined = errors.Join(joined, fmt.Errorf("shutdown http server: %w", err))
	}

	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil {
			joined = errors.Join(joined, fmt.Errorf("get sql db: %w", err))
		} else if err := sqlDB.Close(); err != nil {
			joined = errors.Join(joined, fmt.Errorf("close postgres: %w", err))
		}
	}

	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			joined = errors.Join(joined, fmt.Errorf("close redis: %w", err))
		}
	}

	return joined
}
