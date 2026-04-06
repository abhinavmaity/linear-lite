package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(databaseURL, migrationsPath, direction string, steps int) error {
	m, err := migrate.New("file://"+migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("initialize migrate: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	switch direction {
	case "up":
		if err := runUp(m, steps); err != nil {
			return err
		}
	case "down":
		if err := runDown(m, steps); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported migration direction %q, use up or down", direction)
	}

	return nil
}

func runUp(m *migrate.Migrate, steps int) error {
	if steps > 0 {
		if err := m.Steps(steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("run up steps: %w", err)
		}
		return nil
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run up: %w", err)
	}
	return nil
}

func runDown(m *migrate.Migrate, steps int) error {
	if steps > 0 {
		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("run down steps: %w", err)
		}
		return nil
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run down: %w", err)
	}
	return nil
}
