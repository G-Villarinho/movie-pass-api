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

	var payload domain.CinemaPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	response, err := c.cinemaService.Create(ctx.Request().Context(), payload)
	if err != nil {
		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c *cinemaHandler) GetByID(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "cinema"),
		slog.String("func", "GetByID"),
	)

	param := ctx.Param("cinemaId")
	cinemaID, err := uuid.Parse(param)
	if err != nil {
		log.Warn("Invalid cinema ID provided", slog.String("cinemaId", param), slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid ID", "The provided cinema ID is not a valid UUID.")
	}

	response, err := c.cinemaService.GetByID(ctx.Request().Context(), cinemaID)
	if err != nil {
		if errors.Is(err, domain.ErrCinemaNotFound) {
			return ctx.NoContent(http.StatusNoContent)
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *cinemaHandler) Delete(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "cinema"),
		slog.String("func", "Delete"),
	)

	param := ctx.Param("cinemaId")
	cinemaID, err := uuid.Parse(param)
	if err != nil {
		log.Warn("Invalid cinema ID provided", slog.String("cinemaId", param), slog.String("error", err.Error()))
		return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, nil, "Invalid ID", "The provided cinema ID is not a valid UUID.")
	}

	if err := c.cinemaService.Delete(ctx.Request().Context(), cinemaID); err != nil {
		if errors.Is(err, domain.ErrCinemaNotFound) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "Not Found", "The cinema you are trying to delete does not exist.")
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *cinemaHandler) GetAll(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "cinema"),
		slog.String("func", "GetAll"),
	)

	limit := ctx.QueryParam("limit")
	page := ctx.QueryParam("page")
	sort := ctx.QueryParam("sort")

	pagination := &domain.Pagination{}
	pagination.SetLimit(limit)
	pagination.SetPage(page)
	pagination.SetSort(sort)

	response, err := c.cinemaService.GetAll(ctx.Request().Context(), pagination)
	if err != nil {
		if errors.Is(err, domain.ErrCinemaNotFound) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, nil, "No Cinemas Found", "There are currently no cinemas available in the system. Please try again later or contact support if you believe this is a mistake.")
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}
