package handler

import (
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func SetupRoutes(e *echo.Echo, i *do.Injector) {
	setupUserRoutes(e, i)
}

func setupUserRoutes(e *echo.Echo, i *do.Injector) {
	userHandler, err := do.Invoke[domain.UserHandler](i)
	if err != nil {
		panic(err)
	}

	group := e.Group("/v1/users")
	group.POST("", userHandler.Create)
}
