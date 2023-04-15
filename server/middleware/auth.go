package middleware

import (
	"errors"
	"net/http"

	"remote_monitoring_and_controlling/pkg/jwt"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(auth *jwt.JWT) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(ctx echo.Context, err error) error {
			var errorText string

			if errors.Is(err, jwt.ErrExpired) || errors.Is(err, jwt.ErrUserGet) {
				errorText = err.Error()
			} else {
				errorText = "Failed to validate token"
			}

			return ctx.JSON(http.StatusUnauthorized, map[string]any{
				"error": errorText,
			})
		},
		ParseTokenFunc: func(ctx echo.Context, accessToken string) (any, error) {
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
