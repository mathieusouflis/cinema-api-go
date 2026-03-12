package redisRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenRepository struct {
	client *redis.Client
}

func NewRedisTokenRepository(client *redis.Client) *RedisTokenRepository {
	return &RedisTokenRepository{client: client}
}

func refreshKey(userID string) string {
	return fmt.Sprintf("refresh_token:%s", userID)
}

func (r *RedisTokenRepository) StoreRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error {
	return r.client.Set(ctx, refreshKey(userID), token, ttl).Err()
}

func (r *RedisTokenRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	token, err := r.client.Get(ctx, refreshKey(userID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return token, err
}

func (r *RedisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	return r.client.Del(ctx, refreshKey(userID)).Err()
}
