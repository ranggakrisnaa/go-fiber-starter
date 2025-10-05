package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/ranggakrisnaa/go-fiber-starter/middlewares"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/auth/service"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/user/controller"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/constants"
	"github.com/samber/do"
)

func RegisterRoutes(app *fiber.App, injector *do.Injector) {
	userController := do.MustInvoke[controller.UserController](injector)
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)

	userRoutes := app.Group("/api/user")
	{
		userRoutes.Get("", userController.GetAllUser)
		userRoutes.Get("/me", middlewares.Authenticate(jwtService), userController.Me)
		userRoutes.Put("/:id", middlewares.Authenticate(jwtService), userController.Update)
		userRoutes.Delete("/:id", middlewares.Authenticate(jwtService), userController.Delete)
	}
}
