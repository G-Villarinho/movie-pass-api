package main

import (
	"context"
	"fmt"
	"time"

	"github.com/GSVillas/movie-pass-api/api/handler"
	"github.com/GSVillas/movie-pass-api/config"
	"github.com/GSVillas/movie-pass-api/repository"
	"github.com/GSVillas/movie-pass-api/service"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	e := echo.New()
	i := do.New()

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	do.Provide(i, handler.NewUserHandler)

	do.Provide(i, service.NewUserService)

	do.Provide(i, repository.NewUserRepository)

	handler.SetupRoutes(e, i)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Env.APIPort)))
}
