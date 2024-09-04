package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type movieHandler struct {
	i            *do.Injector
	movieService domain.MovieService
}

func NewMovieHandler(i *do.Injector) (domain.MovieHandler, error) {
	movieService, err := do.Invoke[domain.MovieService](i)
	if err != nil {
		return nil, err
	}

	return &movieHandler{
		i:            i,
		movieService: movieService,
	}, nil
}

func (m *movieHandler) GetAllIndicativeRating(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "movie"),
		slog.String("func", "GetAllIndicativeRating"),
	)

	log.Info("Initializing get all indicative rating process")

	response, err := m.movieService.GetAllIndicativeRating(ctx.Request().Context())
	if err != nil {
		if errors.Is(err, domain.ErrIndicativeRatingsNotFound) {
			log.Warn("No indicative rating available", slog.String("error", err.Error()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "No indicative rating available Found", "There are currently no indicative rating available in the system. Please try again later or contact support if you believe this is a mistake.")
		}

		log.Error("Failed to get all indicative ratings", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("Get all indicative ratings executed successfully")
	return ctx.JSON(http.StatusOK, response)
}
