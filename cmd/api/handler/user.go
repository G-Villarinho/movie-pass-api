package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GSVillas/movie-pass-api/domain"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type userHandler struct {
	i           *do.Injector
	userService domain.UserService
}

func NewUserHandler(i *do.Injector) (domain.UserHandler, error) {
	userService, err := do.Invoke[domain.UserService](i)
	if err != nil {
		return nil, err
	}

	return &userHandler{
		i:           i,
		userService: userService,
	}, nil
}

func (u *userHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "Create"),
	)

	var payload domain.UserPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Failed to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	if err := u.userService.Create(ctx.Request().Context(), payload); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyRegister) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "conflict", "The email already registered. Please try again with a different email.")
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}

func (u *userHandler) SignIn(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "Create"),
	)

	var payload domain.SignInPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Failed to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	response, err := u.userService.SignIn(ctx.Request().Context(), payload)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidPassword) {
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusUnauthorized, nil, "Unauthorized", "Invalid email or password. Please check your credentials and try again.")
		}

		log.Error(err.Error())
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (u *userHandler) CreateAdmin(ctx echo.Context) error {
	panic("unimplemented")
}
