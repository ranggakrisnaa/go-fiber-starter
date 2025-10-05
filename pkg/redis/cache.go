package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type RedisCacheService struct {
	client *redis.Client
}

func NewRedisCacheService(client *redis.Client) CacheService {
	return &RedisCacheService{
		client: client,
	}
}

func (r *RedisCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

func (r *RedisCacheService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCacheService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}