package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
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

func (m *movieHandler) GetAllByUserID(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "movie"),
		slog.String("func", "GetAllByUserID"),
	)

	log.Info("Starting process to retrieve all movies by user ID")

	response, err := m.movieService.GetAllByUserID(ctx.Request().Context())
	if err != nil {
		if errors.Is(err, domain.ErrMoviesNotFoundByUserID) {
			log.Warn("No movies found for the user", slog.String("error", err.Error()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Movies Not Found", "No movies were found for the current user. If you believe this is a mistake, please contact support.")
		}

		log.Error("Failed to retrieve movies by user ID", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("Successfully retrieved movies for user", slog.String("userID", ctx.Request().Header.Get("userID")))
	return ctx.JSON(http.StatusOK, response)
}

func (m *movieHandler) Update(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "movie"),
		slog.String("func", "Update"),
	)

	log.Info("Initializing update movie process")

	movieIDParam := ctx.Param("id")
	movieID, err := uuid.Parse(movieIDParam)
	if err != nil {
		log.Warn("Invalid Movie ID", slog.String("movieID", movieIDParam))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid Movie ID", "The provided movie ID is not a valid UUID.")
	}

	var payload domain.MovieUpdatePayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Failed to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		log.Warn("Validation failed", slog.Any("errors", validationErrors))
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	response, err := m.movieService.Update(ctx.Request().Context(), movieID, payload)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFoundInContext) {
			log.Warn("User not found in session")
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusUnauthorized, nil, "Unauthorized", "User is not authenticated or session has expired.")
		}

		if errors.Is(err, domain.ErrMoviesNotFound) {
			log.Warn("Movie not found", slog.String("movieID", movieID.String()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Movie Not Found", "The movie you are trying to update does not exist.")
		}

		if errors.Is(err, domain.ErrMovieNotBelongUser) {
			log.Warn("User does not own the movie", slog.String("userID", ctx.Request().Header.Get("userID")), slog.String("movieID", movieID.String()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusForbidden, nil, "Forbidden", "You are not allowed to update this movie because it does not belong to you.")
		}

		if errors.Is(err, domain.ErrIndicativeRatingNotFound) {
			log.Warn("Indicative rating not found", slog.String("indicativeRatingID", payload.IndicativeRatingID.String()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid Indicative Rating", "The provided indicative rating does not exist.")
		}

		log.Error("Failed to update movie", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("Movie updated successfully")
	return ctx.JSON(http.StatusOK, response)
}
