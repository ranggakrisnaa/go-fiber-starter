package middlewares

import (
	"github.com/gofiber/fiber/v3"
)

func CORSMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Set("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Method() == "OPTIONS" {
			return c.Status(fiber.StatusNoContent).Send(nil)
		}

		return c.Next()
	}
}
