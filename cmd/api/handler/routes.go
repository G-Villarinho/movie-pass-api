package handler

import (
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/GSVillas/movie-pass-api/middleware"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func SetupRoutes(e *echo.Echo, i *do.Injector) {
	setupUserRoutes(e, i)
	setupCinemaRoutes(e, i)
	setupMovieRoutes(e, i)
}

func setupUserRoutes(e *echo.Echo, i *do.Injector) {
	userHandler, err := do.Invoke[domain.UserHandler](i)
	if err != nil {
		panic(err)
	}

	group := e.Group("/v1/users")
	group.POST("", userHandler.Create)
	group.POST("/admin", userHandler.Create, middleware.RequireAdminLevel(domain.AdminRoleLevel3))
	group.POST("/sign-in", userHandler.SignIn)
}

func setupCinemaRoutes(e *echo.Echo, i *do.Injector) {
	cinemaHandler, err := do.Invoke[domain.CinemaHandler](i)
	if err != nil {
		panic(err)
	}

	group := e.Group("/v1/cinemas", middleware.EnsureAuthenticated(i), middleware.RequireAdminLevel(domain.AdminRoleLevel1))
	group.POST("", cinemaHandler.Create)
	group.GET("", cinemaHandler.GetAll)
	group.GET("/:id", cinemaHandler.GetByID)
	group.DELETE("/:id", cinemaHandler.Delete)
}

func setupMovieRoutes(e *echo.Echo, i *do.Injector) {
	movieHandler, err := do.Invoke[domain.MovieHandler](i)
	if err != nil {
		panic(err)
	}

	group := e.Group("/v1/movies")
	group.GET("/indicative-rating", movieHandler.GetAllIndicativeRatings)
	group.POST("", movieHandler.Create, middleware.EnsureAuthenticated(i), middleware.RequireAdminLevel(domain.AdminRoleLevel1))
	group.GET("", movieHandler.GetAllByUserID, middleware.EnsureAuthenticated(i), middleware.RequireAdminLevel(domain.AdminRoleLevel1))
	group.PUT("/:id", movieHandler.Update, middleware.EnsureAuthenticated(i), middleware.RequireAdminLevel(domain.AdminRoleLevel1))
	group.PUT("/:id", movieHandler.Update, middleware.EnsureAuthenticated(i), middleware.RequireAdminLevel(domain.AdminRoleLevel1))

}
