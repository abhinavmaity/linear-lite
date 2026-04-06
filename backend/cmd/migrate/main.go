package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abhinavmaity/linear-lite/backend/internal/database"
)

func main() {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		slog.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	migrationsPath := strings.TrimSpace(os.Getenv("MIGRATIONS_PATH"))
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	direction := strings.ToLower(strings.TrimSpace(os.Getenv("MIGRATION_DIRECTION")))
	if direction == "" {
		direction = "up"
	}

	steps, err := parseSteps(os.Getenv("MIGRATION_STEPS"))
	if err != nil {
		slog.Error("invalid migration steps", "error", err)
		os.Exit(1)
	}

	absMigrationsPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		slog.Error("failed to resolve migrations path", "error", err)
		os.Exit(1)
	}

	slog.Info("running migrations", "direction", direction, "steps", steps, "path", absMigrationsPath)

	if err := database.Migrate(databaseURL, absMigrationsPath, direction, steps); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("migration completed")
}

func parseSteps(raw string) (int, error) {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return 0, nil
	}

	steps, err := strconv.Atoi(clean)
	if err != nil {
		return 0, fmt.Errorf("MIGRATION_STEPS must be an integer: %w", err)
	}
	if steps < 0 {
		return 0, fmt.Errorf("MIGRATION_STEPS must be >= 0")
	}

	return steps, nil
}
