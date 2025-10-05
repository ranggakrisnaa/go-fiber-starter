package logs

import (
	"github.com/gofiber/fiber/v3"
	"github.com/samber/do"
)

func RegisterRoutes(app *fiber.App, injector *do.Injector) {
	logController := &LogController{}

	app.Get("/logs", logController.GetLogs)
	app.Get("/logs/:month", logController.GetLogs)
}
