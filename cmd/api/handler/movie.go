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

	response, err := m.movieService.GetAllIndicativeRating(ctx.Request().Context())
	if err != nil {
		if errors.Is(err, domain.ErrIndicativeRatingsNotFound) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "No indicative rating available Found", "There are currently no indicative rating available in the system. Please try again later or contact support if you believe this is a mistake.")
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (m *movieHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "movie"),
		slog.String("func", "Create"),
	)

	form, err := ctx.MultipartForm()
	if err != nil {
		log.Error("Failed to parse multipart form", slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Form Parsing Error", "Failed to parse multipart form data.")
	}

	duration, err := strconv.Atoi(ctx.FormValue("duration"))
	if err != nil {
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid Duration", "The provided duration must be a valid positive number.")
	}

	indicativeRatingID, err := uuid.Parse(ctx.FormValue("indicativeRatingId"))
	if err != nil {
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid IndicativeRatingID", "The provided indicative rating ID is not a valid UUID.")
	}

	payload := domain.MoviePayload{
		Title:              ctx.FormValue("title"),
		Duration:           duration,
		IndicativeRatingID: indicativeRatingID,
		Images:             form.File["images"],
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	response, err := m.movieService.Create(ctx.Request().Context(), payload)
	if err != nil {
		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (m *movieHandler) GetAllByUserID(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "movie"),
		slog.String("func", "GetAllByUserID"),
	)

	response, err := m.movieService.GetAllByUserID(ctx.Request().Context())
	if err != nil {
		if errors.Is(err, domain.ErrMoviesNotFoundByUserID) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Movies Not Found", "No movies were found for the current user. If you believe this is a mistake, please contact support.")
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (m *movieHandler) Update(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "movie"),
		slog.String("func", "Update"),
	)

	movieIDParam := ctx.Param("id")
	movieID, err := uuid.Parse(movieIDParam)
	if err != nil {
		log.Warn("Invalid Movie ID", slog.String("movieID", movieIDParam))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid Movie ID", "The provided movie ID is not a valid UUID.")
	}

	var payload domain.MovieUpdatePayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	response, err := m.movieService.Update(ctx.Request().Context(), movieID, payload)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFoundInContext) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusUnauthorized, nil, "Unauthorized", "User is not authenticated or session has expired.")
		}

		if errors.Is(err, domain.ErrMoviesNotFound) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Movie Not Found", "The movie you are trying to update does not exist.")
		}

		if errors.Is(err, domain.ErrMovieNotBelongUser) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusForbidden, nil, "Forbidden", "You are not allowed to update this movie because it does not belong to you.")
		}

		if errors.Is(err, domain.ErrIndicativeRatingNotFound) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid Indicative Rating", "The provided indicative rating does not exist.")
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (m *movieHandler) Delete(ctx echo.Context) error {
	panic("unimplemented")
}
