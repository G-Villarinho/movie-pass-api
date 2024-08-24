package middleware

import (
	"context"
	"strings"

	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func EnsureAuthenticated(i *do.Injector) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			sessionService, err := do.Invoke[domain.SessionService](i)
			if err != nil {
				return domain.InternalServerAPIErrorResponse(ctx)
			}

			authorizationHeader := ctx.Request().Header.Get("Authorization")
			if authorizationHeader == "" {
				return domain.AccessDeniedAPIErrorResponse(ctx)
			}

			content := strings.Split(authorizationHeader, " ")
			if len(content) != 2 {
				return domain.AccessDeniedAPIErrorResponse(ctx)
			}

			session, err := sessionService.GetSession(ctx.Request().Context(), content[1])
			if err != nil {
				if err == domain.ErrTokenInvalid || err == domain.ErrSessionMismatch || err == domain.ErrSessionNotFound {
					return domain.AccessDeniedAPIErrorResponse(ctx)
				}

				return domain.InternalServerAPIErrorResponse(ctx)
			}

			newCtx := context.WithValue(ctx.Request().Context(), domain.SessionKey, session)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))

			return next(ctx)
		}
	}
}
