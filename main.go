package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GSVillas/movie-pass-api/api/handler"
	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/config/database"
	"github.com/GSVillas/movie-pass-api/repository"
	"github.com/GSVillas/movie-pass-api/service"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	e := echo.New()
	i := do.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("Fail to connect to mysql: ", err)
	}

	redisClient, err := database.NewRedisConnection(ctx)
	if err != nil {
		log.Fatal("Fail to connect to redis: ", err)
	}

	do.Provide(i, func(i *do.Injector) (*gorm.DB, error) {
		return db, nil
	})

	do.Provide(i, func(i *do.Injector) (*redis.Client, error) {
		return redisClient, nil
	})

	do.Provide(i, handler.NewUserHandler)

	do.Provide(i, service.NewUserService)
	do.Provide(i, service.NewSessionService)

	do.Provide(i, repository.NewUserRepository)
	do.Provide(i, repository.NewSessionRepository)

	handler.SetupRoutes(e, i)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Env.APIPort)))
}
