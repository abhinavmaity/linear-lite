package validation

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Pagination struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func ParseUUIDParam(c *gin.Context, key string) (uuid.UUID, *apperrors.AppError) {
	value := strings.TrimSpace(c.Param(key))
	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, apperrors.Validation("invalid path parameter", apperrors.FieldErrors{
			key: "must be a valid UUID",
		})
	}
	return parsed, nil
}

func ParsePagination(c *gin.Context, defaultLimit, maxLimit int) (Pagination, *apperrors.AppError) {
	page, err := parsePositiveInt(c.DefaultQuery("page", "1"))
	if err != nil {
		return Pagination{}, apperrors.Validation("invalid pagination parameters", apperrors.FieldErrors{
			"page": "must be a positive integer",
		})
	}

	limit, err := parsePositiveInt(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if err != nil {
		return Pagination{}, apperrors.Validation("invalid pagination parameters", apperrors.FieldErrors{
			"limit": "must be a positive integer",
		})
	}
	if limit > maxLimit {
		return Pagination{}, apperrors.Validation("invalid pagination parameters", apperrors.FieldErrors{
			"limit": "must be less than or equal to " + strconv.Itoa(maxLimit),
		})
	}

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}, nil
}

func ParseSortOrder(raw string) (string, *apperrors.AppError) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return "desc", nil
	}
	if value != "asc" && value != "desc" {
		return "", apperrors.Validation("invalid sort order", apperrors.FieldErrors{
			"order": "must be asc or desc",
		})
	}
	return value, nil
}

func ParseSortField(raw, defaultValue string, allowed []string) (string, *apperrors.AppError) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return defaultValue, nil
	}

	for _, option := range allowed {
		if value == option {
			return value, nil
		}
	}

	return "", apperrors.Validation("invalid sort field", apperrors.FieldErrors{
		"sort_by": "must be one of: " + strings.Join(allowed, ", "),
	})
}

func ValidateEnum(field string, raw string, allowed []string) *apperrors.AppError {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}

	for _, option := range allowed {
		if value == option {
			return nil
		}
	}

	return apperrors.Validation("invalid enum value", apperrors.FieldErrors{
		field: "must be one of: " + strings.Join(allowed, ", "),
	})
}

func ParseDate(field, raw string) (time.Time, *apperrors.AppError) {
	value := strings.TrimSpace(raw)
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, apperrors.Validation("invalid date", apperrors.FieldErrors{
			field: "must use YYYY-MM-DD format",
		})
	}
	return parsed, nil
}

func NormalizeOptionalString(raw string) *string {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return nil
	}
	return &clean
}

func ParseRepeatedQuery(c *gin.Context, key string) []string {
	values := c.QueryArray(key)
	out := make([]string, 0, len(values))
	for _, value := range values {
		clean := strings.TrimSpace(value)
		if clean != "" {
			out = append(out, clean)
		}
	}
	return out
}

func ParseOptionalUUIDQuery(c *gin.Context, key string) (*uuid.UUID, *apperrors.AppError) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil, apperrors.Validation("invalid query parameter", apperrors.FieldErrors{
			key: "must be a valid UUID",
		})
	}

	return &parsed, nil
}

func ParseOptionalBoolQuery(c *gin.Context, key string, defaultValue bool) (bool, *apperrors.AppError) {
	rawValues, exists := c.GetQueryArray(key)
	if !exists || len(rawValues) == 0 {
		return defaultValue, nil
	}

	value := strings.ToLower(strings.TrimSpace(rawValues[len(rawValues)-1]))
	if value == "true" || value == "1" {
		return true, nil
	}
	if value == "false" || value == "0" {
		return false, nil
	}

	return false, apperrors.Validation("invalid query parameter", apperrors.FieldErrors{
		key: "must be a boolean",
	})
}

func ParseDistinctUUIDArray(field string, raw []string) ([]uuid.UUID, *apperrors.AppError) {
	if len(raw) == 0 {
		return nil, nil
	}

	out := make([]uuid.UUID, 0, len(raw))
	seen := make(map[uuid.UUID]struct{}, len(raw))
	for idx, value := range raw {
		clean := strings.TrimSpace(value)
		parsed, err := uuid.Parse(clean)
		if err != nil {
			return nil, apperrors.Validation("invalid query parameter", apperrors.FieldErrors{
				field: fmt.Sprintf("item %d must be a valid UUID", idx+1),
			})
		}
		if _, ok := seen[parsed]; ok {
			return nil, apperrors.Validation("invalid query parameter", apperrors.FieldErrors{
				field: "must not contain duplicates",
			})
		}
		seen[parsed] = struct{}{}
		out = append(out, parsed)
	}

	return out, nil
}

func QueryHasKey(c *gin.Context, key string) bool {
	parsed, err := url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		return false
	}
	_, ok := parsed[key]
	return ok
}

func parsePositiveInt(raw string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, err
	}
	if value <= 0 {
		return 0, strconv.ErrSyntax
	}
	return value, nil
}
