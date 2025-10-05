package providers

import (
	"github.com/ranggakrisnaa/go-fiber-starter/config"
	authController "github.com/ranggakrisnaa/go-fiber-starter/modules/auth/controller"
	authRepo "github.com/ranggakrisnaa/go-fiber-starter/modules/auth/repository"
	authService "github.com/ranggakrisnaa/go-fiber-starter/modules/auth/service"
	userController "github.com/ranggakrisnaa/go-fiber-starter/modules/user/controller"
	"github.com/ranggakrisnaa/go-fiber-starter/modules/user/repository"
	userService "github.com/ranggakrisnaa/go-fiber-starter/modules/user/service"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/constants"
	"github.com/ranggakrisnaa/go-fiber-starter/pkg/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func InitDatabase(injector *do.Injector) {
	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return config.SetUpDatabaseConnection(), nil
	})
}

func InitKafka(injector *do.Injector) {
	do.ProvideNamed(injector, constants.KafkaProducer, func(i *do.Injector) (kafka.KafkaProducer, error) {
		return kafka.NewKafkaProducer()
	})

	do.ProvideNamed(injector, constants.KafkaConsumer, func(i *do.Injector) (kafka.KafkaConsumer, error) {
		return kafka.NewKafkaConsumer()
	})
}

func InitRedis(injector *do.Injector) {
	do.ProvideNamed(injector, constants.RedisClient, func(i *do.Injector) (*redis.Client, error) {
		return config.NewRedisConfig().GetRedisClient(), nil
	})
}

func RegisterDependencies(injector *do.Injector) {
	InitDatabase(injector)
	InitKafka(injector)
	InitRedis(injector)

	do.ProvideNamed(injector, constants.JWTService, func(i *do.Injector) (authService.JWTService, error) {
		return authService.NewJWTService(), nil
	})

	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := authRepo.NewRefreshTokenRepository(db)

	userService := userService.NewUserService(userRepository, db)
	authService := authService.NewAuthService(userRepository, refreshTokenRepository, jwtService, db)

	do.Provide(
		injector, func(i *do.Injector) (userController.UserController, error) {
			return userController.NewUserController(i, userService), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (authController.AuthController, error) {
			return authController.NewAuthController(i, authService), nil
		},
	)
}
