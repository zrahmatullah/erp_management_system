package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

func (r *SessionRepository) BlacklistToken(ctx context.Context, token string, expiry time.Duration) error {
	return r.client.Set(ctx, "bl:"+token, "true", expiry).Err()
}

func (r *SessionRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	res, err := r.client.Get(ctx, "bl:"+token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return res == "true", nil
}

func (r *SessionRepository) SetCache(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiry).Err()
}

func (r *SessionRepository) GetCache(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *SessionRepository) DeleteCache(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
