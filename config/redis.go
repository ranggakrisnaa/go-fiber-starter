package config

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisConfig() *RedisConfig {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Failed to load .env file: " + err.Error())
	}

	return &RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // Default DB
	}
}

func (r *RedisConfig) GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Host + ":" + r.Port,
		Password: r.Password,
		DB:       r.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return client
}
