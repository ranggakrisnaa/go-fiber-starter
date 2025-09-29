package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/auth/controller"
	"github.com/samber/do"
)

func RegisterRoutes(app *fiber.App, injector *do.Injector) {
	authController := do.MustInvoke[controller.AuthController](injector)

	authRoutes := app.Group("/api/auth")
	{
		authRoutes.Post("/register", authController.Register)
		authRoutes.Post("/login", authController.Login)
		authRoutes.Post("/refresh", authController.RefreshToken)
		authRoutes.Post("/logout", authController.Logout)
		authRoutes.Post("/send-verification-email", authController.SendVerificationEmail)
		authRoutes.Post("/verify-email", authController.VerifyEmail)
		authRoutes.Post("/send-password-reset", authController.SendPasswordReset)
		authRoutes.Post("/reset-password", authController.ResetPassword)
	}
}
