package middlewares

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/auth/service"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/user/dto"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/utils"
)

func RequireRole(jwtService service.JWTService, allowedRoles ...string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		// Ambil token dari context (harus sudah diset oleh Authenticate middleware)
		token := ctx.Locals("token")
		if token == nil {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_FOUND, nil)
			return ctx.Status(fiber.StatusUnauthorized).JSON(response)
		}

		tokenString, ok := token.(string)
		if !ok {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			return ctx.Status(fiber.StatusUnauthorized).JSON(response)
		}

		userRoles, err := jwtService.GetRolesByToken(tokenString)
		if err != nil {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			return ctx.Status(fiber.StatusUnauthorized).JSON(response)
		}

		hasPermission := false
		for _, userRole := range userRoles {
			if slices.Contains(allowedRoles, userRole) {
				hasPermission = true
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, "Access denied: insufficient permissions", nil)
			return ctx.Status(fiber.StatusForbidden).JSON(response)
		}

		ctx.Locals("user_roles", userRoles)
		return ctx.Next()
	}
}

func RequireAdmin(jwtService service.JWTService) fiber.Handler {
	return RequireRole(jwtService, "admin")
}

func RequireUser(jwtService service.JWTService) fiber.Handler {
	return RequireRole(jwtService, "user", "admin")
}
