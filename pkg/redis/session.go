package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID    string      `json:"user_id"`
	Username  string      `json:"username"`
	Role      string      `json:"role"`
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

type SessionService interface {
	CreateSession(ctx context.Context, sessionID string, data SessionData, expiration time.Duration) error
	GetSession(ctx context.Context, sessionID string) (*SessionData, error)
	DeleteSession(ctx context.Context, sessionID string) error
	UpdateSession(ctx context.Context, sessionID string, data SessionData) error
}

type RedisSessionService struct {
	client *redis.Client
}

func NewRedisSessionService(client *redis.Client) SessionService {
	return &RedisSessionService{
		client: client,
	}
}

func (r *RedisSessionService) CreateSession(ctx context.Context, sessionID string, data SessionData, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "session:"+sessionID, jsonData, expiration).Err()
}

func (r *RedisSessionService) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	val, err := r.client.Get(ctx, "session:"+sessionID).Result()
	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	if err := json.Unmarshal([]byte(val), &sessionData); err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (r *RedisSessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, "session:"+sessionID).Err()
}

func (r *RedisSessionService) UpdateSession(ctx context.Context, sessionID string, data SessionData) error {
	expiresIn := time.Until(data.ExpiresAt)
	if expiresIn <= 0 {
		return r.DeleteSession(ctx, sessionID)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "session:"+sessionID, jsonData, expiresIn).Err()
}
