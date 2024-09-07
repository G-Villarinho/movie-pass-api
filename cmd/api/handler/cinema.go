package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type cinemaHandler struct {
	i             *do.Injector
	cinemaService domain.CinemaService
}

func NewCinemaHandler(i *do.Injector) (domain.CinemaHandler, error) {
	cinemaService, err := do.Invoke[domain.CinemaService](i)
	if err != nil {
		return nil, err
	}

	return &cinemaHandler{
		i:             i,
		cinemaService: cinemaService,
	}, nil

}

func (c *cinemaHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "cinema"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing cinema creation process")

	var payload domain.CinemaPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Failed to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		log.Warn("Validation failed", slog.Any("errors", validationErrors))
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	response, err := c.cinemaService.Create(ctx.Request().Context(), payload)
	if err != nil {
		log.Error("Fail to create user", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("Cinema created successfully")
	return ctx.JSON(http.StatusCreated, response)
}

func (c *cinemaHandler) GetByID(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "cinema"),
		slog.String("func", "GetByID"),
	)

	log.Info("Initializing cinema get by id process")

	param := ctx.Param("cinemaId")

	cinemaID, err := uuid.Parse(param)
	if err != nil {
		log.Warn("Invalid cinema ID provided", slog.String("cinemaId", param), slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid ID", "The provided cinema ID is not a valid UUID.")
	}

	response, err := c.cinemaService.GetByID(ctx.Request().Context(), cinemaID)
	if err != nil {
		if errors.Is(err, domain.ErrCinemaNotFound) {
			log.Warn("No cinema found", slog.String("error", err.Error()))
			return ctx.NoContent(http.StatusNoContent)
		}

		log.Error("Fail to create user", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("Get cinema by id executed successfully")
	return ctx.JSON(http.StatusOK, response)
}

func (c *cinemaHandler) Delete(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "cinema"),
		slog.String("func", "Delete"),
	)

	log.Info("Initializing cinema delete process")

	param := ctx.Param("cinemaId")

	cinemaID, err := uuid.Parse(param)
	if err != nil {
		log.Warn("Invalid cinema ID provided", slog.String("cinemaId", param), slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid ID", "The provided cinema ID is not a valid UUID.")
	}

	if err := c.cinemaService.Delete(ctx.Request().Context(), cinemaID); err != nil {
		if errors.Is(err, domain.ErrCinemaNotFound) {
			log.Warn("Cinema not found for deletion", slog.String("cinemaId", cinemaID.String()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The cinema you are trying to delete does not exist.")
		}

		log.Error("Failed to delete cinema", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("Cinema deleted successfully", slog.String("cinemaId", cinemaID.String()))
	return ctx.NoContent(http.StatusOK)
}

func (c *cinemaHandler) GetAll(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "cinema"),
		slog.String("func", "GetAll"),
	)

	log.Info("Initializing get all cinemas process")

	response, err := c.cinemaService.GetAll(ctx.Request().Context())
	if err != nil {
		if errors.Is(err, domain.ErrCinemaNotFound) {
			log.Warn("No cinemas available", slog.String("error", err.Error()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "No Cinemas Found", "There are currently no cinemas available in the system. Please try again later or contact support if you believe this is a mistake.")
		}

		log.Error("Failed to get all cinemas", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("Get all cinemas executed successfully")
	return ctx.JSON(http.StatusOK, response)
}
