package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/constants"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/redis"
	"github.com/samber/do"
)

func RateLimiter(injector *do.Injector, limit int64, window time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		rateLimiter := redis.NewRedisRateLimiterService(
			do.MustInvokeNamed[*redis.Client](injector, constants.RedisClient),
		)

		clientIP := c.IP()
		path := c.Path()
		key := "rate_limit:" + clientIP + ":" + path

		allowed, err := rateLimiter.AllowRequest(c.Context(), key, limit, window)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}

		return c.Next()
	}
}
