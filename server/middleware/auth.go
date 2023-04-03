package middleware

import (
	"net/http"
	"remote_monitoring_and_controlling/pkg/jwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AuthMiddleware(auth *jwt.JWT) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		ErrorHandlerWithContext: func(err error, ctx echo.Context) error {
			errorText := ""

			switch err {
			case jwt.ErrExpired, jwt.ErrUserGet:
				errorText = err.Error()
			default:
				errorText = "Failed to validate token"
			}

			return ctx.JSON(http.StatusUnauthorized, map[string]any{
				"error": errorText,
			})
		},
		ParseTokenFunc: func(accessToken string, ctx echo.Context) (any, error) {
			userID, err := auth.ValidateToken(ctx.Request().Context(), accessToken, true)
			if err != nil {
				return nil, err
			}

			ctx.Set("token", accessToken)
			ctx.Set("userID", userID)
			return userID, nil
		},
	})
}
