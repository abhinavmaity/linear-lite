package services

import (
	"strconv"
	"strings"

	cachepkg "github.com/abhinavmaity/linear-lite/backend/internal/cache"
)

func buildListCacheKey(prefix string, parts ...string) string {
	return prefix + ":list:" + cachepkg.HashParts(parts...)
}

func buildDetailCacheKey(prefix, id string) string {
	return prefix + ":detail:" + strings.TrimSpace(id)
}

func intToCachePart(value int) string {
	return strconv.Itoa(value)
}
