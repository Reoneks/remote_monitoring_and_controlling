package server

import (
	"context"
	"net/http"
	"project/config"
	"project/pkg/jwt"
	"project/server/handlers"
	"project/server/middleware"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type HTTPServer struct {
	router *echo.Echo

	handlers *handlers.Handler
	auth     *jwt.JWT

	appAddr string
}

func NewHTTPServer(
	cfg *config.Config,
	handlers *handlers.Handler,
	auth *jwt.JWT,
) *HTTPServer {
	return &HTTPServer{
		appAddr:  cfg.AppAddr,
		handlers: handlers,
		auth:     auth,
	}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	s.router = echo.New()
	s.router.Use(middleware.LoggerMiddleware(), middleware.CorsMiddleware(), middleware.RecoverMiddleware())

	s.router.POST("/login", s.handlers.Login)
	s.router.POST("/2fa", s.handlers.TwoFA)
	s.router.POST("/register", s.handlers.Register)

	private := s.router.Group("", middleware.AuthMiddleware(s.auth))
	private.POST("/logout", s.handlers.Logout)
	private.PATCH("/enable_2fa", s.handlers.EnableTwoFA)

	go func() {
		if err := s.router.Start(s.appAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Str("function", "Start").Err(err).Msg("Server start error")
		}
	}()

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
