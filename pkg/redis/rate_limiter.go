package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client = redis.Client

type RateLimiterService interface {
	AllowRequest(ctx context.Context, key string, limit int64, window time.Duration) (bool, error)
	GetRemainingRequests(ctx context.Context, key string, window time.Duration) (int64, error)
}

type RedisRateLimiterService struct {
	client *redis.Client
}

func NewRedisRateLimiterService(client *redis.Client) RateLimiterService {
	return &RedisRateLimiterService{
		client: client,
	}
}

func (r *RedisRateLimiterService) AllowRequest(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	now := time.Now().Unix()
	_ = int64(window.Seconds())

	pipe := r.client.Pipeline()

	incrCmd := pipe.Incr(ctx, key)
	_ = pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	currentCount := incrCmd.Val()

	if currentCount == 1 {
		pipe.Set(ctx, key+":timestamp", now, window)
		pipe.Exec(ctx)
	}

	return currentCount <= limit, nil
}

func (r *RedisRateLimiterService) GetRemainingRequests(ctx context.Context, key string, window time.Duration) (int64, error) {
	currentCount, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return int64(window.Seconds()), nil
	}
	if err != nil {
		return 0, err
	}

	return currentCount, nil
}
