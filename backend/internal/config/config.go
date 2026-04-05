package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	JWTTTL      time.Duration
	CORSOrigins []string
	LogLevel    string
	BcryptCost  int
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:      getString("APP_ENV", "development"),
		Port:        getString("PORT", "8080"),
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		RedisURL:    strings.TrimSpace(os.Getenv("REDIS_URL")),
		JWTSecret:   strings.TrimSpace(os.Getenv("JWT_SECRET")),
		LogLevel:    getString("LOG_LEVEL", "info"),
	}

	jwtTTL, err := getDuration("JWT_TTL", 24*time.Hour)
	if err != nil {
		return Config{}, err
	}
	cfg.JWTTTL = jwtTTL

	bcryptCost, err := getInt("BCRYPT_COST", 12)
	if err != nil {
		return Config{}, err
	}
	cfg.BcryptCost = bcryptCost

	cfg.CORSOrigins = parseCSV(os.Getenv("CORS_ORIGINS"))

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var issues []string

	if c.DatabaseURL == "" {
		issues = append(issues, "DATABASE_URL is required")
	}
	if c.RedisURL == "" {
		issues = append(issues, "REDIS_URL is required")
	}
	if c.JWTSecret == "" {
		issues = append(issues, "JWT_SECRET is required")
	}
	if c.AppEnv != "development" && len(c.JWTSecret) < 32 {
		issues = append(issues, "JWT_SECRET must be at least 32 characters outside development")
	}
	if c.JWTTTL <= 0 {
		issues = append(issues, "JWT_TTL must be greater than 0")
	}
	if c.BcryptCost < 10 || c.BcryptCost > 14 {
		issues = append(issues, "BCRYPT_COST must be between 10 and 14")
	}
	if len(c.CORSOrigins) == 0 {
		issues = append(issues, "CORS_ORIGINS is required")
	}
	if strings.TrimSpace(c.Port) == "" {
		issues = append(issues, "PORT must not be empty")
	}

	if len(issues) > 0 {
		return errors.New(strings.Join(issues, "; "))
	}

	return nil
}

func (c Config) ListenAddr() string {
	if strings.HasPrefix(c.Port, ":") {
		return c.Port
	}
	return ":" + c.Port
}

func getString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getInt(key string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return parsed, nil
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", key, err)
	}
	return parsed, nil
}

func parseCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))

	for _, part := range parts {
		clean := strings.TrimSpace(part)
		if clean != "" {
			out = append(out, clean)
		}
	}

	return out
}
