package middleware

import (
	"github.com/GSVillas/movie-pass-api/domain"
	"github.com/labstack/echo/v4"
)

func RequireAdminLevel(minLevel domain.RoleType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			session, ok := ctx.Get(string(domain.SessionKey)).(*domain.Session)
			if !ok {
				return domain.AccessDeniedAPIErrorResponse(ctx)
			}

			userLevel, userExists := domain.RoleLevels[session.Role]
			minRequiredLevel, minLevelExists := domain.RoleLevels[minLevel]

			if !userExists || !minLevelExists {
				return domain.AccessDeniedAPIErrorResponse(ctx)
			}

			if userLevel < minRequiredLevel {
				return domain.AccessDeniedAPIErrorResponse(ctx)
			}

			return next(ctx)
		}
	}
}
