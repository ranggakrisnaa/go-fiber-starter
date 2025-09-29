package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ranggakrisnaa/go-fiber-starter/middlewares"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/auth"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/user"
	"github.com/ranggakrisnaa/go-fiber-starter/providers"
	"github.com/ranggakrisnaa/go-fiber-starter/script"
	"github.com/samber/do"

	"github.com/common-nighthawk/go-figure"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func args(injector *do.Injector) bool {
	if len(os.Args) > 1 {
		flag := script.Commands(injector)
		return flag
	}

	return true
}

func run(app *fiber.App, injector *do.Injector) {
	app.Use("/assets", static.New("./assets"))

	port := os.Getenv("GOLANG_PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "0.0.0.0:" + port
	} else {
		serve = ":" + port
	}

	myFigure := figure.NewColorFigure("Rangga", "", "green", true)
	myFigure.Print()

	go func() {
		if err := app.Listen(serve); err != nil {
			log.Fatalf("error running server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if err := do.Shutdown[do.Injector](injector); err != nil {
		log.Printf("Error shutting down dependencies: %v", err)
	}

	log.Println("Server exited properly")
}

func main() {
	var (
		injector = do.New()
	)

	providers.RegisterDependencies(injector)

	if !args(injector) {
		return
	}

	app := fiber.New()
	app.Use(middlewares.CORSMiddleware())
	app.Use(middlewares.RateLimiter(injector, 100, time.Second*60))

	// Register module routes
	user.RegisterRoutes(app, injector)
	auth.RegisterRoutes(app, injector)

	run(app, injector)
}
