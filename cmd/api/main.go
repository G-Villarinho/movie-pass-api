package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GSVillas/movie-pass-api/client"
	"github.com/GSVillas/movie-pass-api/cmd/api/handler"
	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/config/database"
	"github.com/GSVillas/movie-pass-api/repository"
	"github.com/GSVillas/movie-pass-api/service"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	e := echo.New()
	i := do.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","method":"${method}","uri":"${uri}","status":${status},"latency":"${latency_human}"}\n`,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.Env.FrontURL},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

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

	do.Provide(i, client.NewCloudFlareService)

	do.Provide(i, handler.NewCinemaHandler)
	do.Provide(i, handler.NewMovieHandler)
	do.Provide(i, handler.NewUserHandler)

	do.Provide(i, service.NewCinemaSevice)
	do.Provide(i, service.NewMovieService)
	do.Provide(i, service.NewUserService)
	do.Provide(i, service.NewSessionService)

	do.Provide(i, repository.NewCinemaRepository)
	do.Provide(i, repository.NewMovieRepository)
	do.Provide(i, repository.NewUserRepository)
	do.Provide(i, repository.NewSessionRepository)

	handler.SetupRoutes(e, i)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Env.APIPort)))
}
