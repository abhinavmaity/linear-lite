package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func HashParts(parts ...string) string {
	joined := strings.Join(parts, "\x1f")
	sum := sha256.Sum256([]byte(joined))
	encoded := hex.EncodeToString(sum[:])
	if len(encoded) > 16 {
		return encoded[:16]
	}
	return encoded
}

func (s *Store) GetJSON(ctx context.Context, key string, out any) (bool, error) {
	value, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if err := json.Unmarshal(value, out); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, payload, ttl).Err()
}

func (s *Store) DeleteByPrefix(ctx context.Context, prefix string) error {
	pattern := prefix
	if !strings.HasSuffix(pattern, "*") {
		pattern += "*"
	}

	var cursor uint64
	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, pattern, 200).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := s.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
