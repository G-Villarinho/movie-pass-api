package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/google/uuid"
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

func (m *movieHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "movie"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing movie creation process")

	form, err := ctx.MultipartForm()
	if err != nil {
		log.Error("Failed to parse multipart form", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Form Parsing Error", "Failed to parse multipart form data.")
	}

	duration, err := strconv.Atoi(ctx.FormValue("duration"))
	if err != nil {
		log.Warn("Invalid duration")
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid Duration", "The provided duration must be a valid positive number.")
	}

	indicativeRatingID, err := uuid.Parse(ctx.FormValue("indicativeRatingId"))
	if err != nil {
		log.Warn("Invalid IndicativeRatingID", slog.String("indicativeRatingId", ctx.FormValue("indicativeRatingId")))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid IndicativeRatingID", "The provided indicative rating ID is not a valid UUID.")
	}

	payload := domain.MoviePayload{
		Title:              ctx.FormValue("title"),
		Duration:           duration,
		IndicativeRatingID: indicativeRatingID,
		Images:             form.File["images"],
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		log.Warn("Validation failed", slog.Any("errors", validationErrors))
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	response, err := m.movieService.Create(ctx.Request().Context(), payload)
	if err != nil {
		log.Error("Failed to create movie", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("Movie created successfully")
	return ctx.JSON(http.StatusCreated, response)
}
