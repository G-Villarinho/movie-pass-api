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

	log.Info("Initializing user creation process")

	var payload domain.UserPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Failed to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		log.Warn("Validation failed", slog.Any("errors", validationErrors))
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	if err := u.userService.Create(ctx.Request().Context(), payload); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyRegister) {
			log.Warn("Fail to create user", slog.String("error", err.Error()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, nil, "conflict", "The email already registered. Please try again with a different email.")
		}

		log.Error("Fail to create user", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("User created successfully")
	return ctx.NoContent(http.StatusCreated)
}

func (u *userHandler) SignIn(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")

	var payload domain.SignInPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Failed to decode JSON payload", slog.String("error", err.Error()))
		return domain.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if validationErrors := payload.Validate(); validationErrors != nil {
		log.Warn("Validation failed", slog.Any("errors", validationErrors))
		return domain.NewValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, validationErrors)
	}

	response, err := u.userService.SignIn(ctx.Request().Context(), payload)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidPassword) {
			log.Warn("Fail to excute user sign in", slog.String("error", err.Error()))
			return domain.NewCustomValidationAPIErrorResponse(ctx, http.StatusUnauthorized, nil, "Unauthorized credentials", "Unauthorized credentials. Review the data sent.")
		}

		log.Error("Fail to create user", slog.String("error", err.Error()))
		return domain.InternalServerAPIErrorResponse(ctx)
	}

	log.Info("user sign in executed succefully")
	return ctx.JSON(http.StatusOK, response)
}
