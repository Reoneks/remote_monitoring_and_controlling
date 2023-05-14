package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

func CheckIPMiddleware(permittedIPs []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if len(permittedIPs) > 0 && !slices.Contains(permittedIPs, c.RealIP()) {
				log.Error().Str("function", "CheckIPMiddleware").Msg(fmt.Sprintf("IP address %s not allowed", c.RealIP()))
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("IP address %s not allowed", c.RealIP()))
			}

			return next(c)
		}
	}
}
