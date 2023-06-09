package server

import (
	"context"
	"errors"
	"net/http"

	"remote_monitoring_and_controlling/config"
	"remote_monitoring_and_controlling/pkg/jwt"
	"remote_monitoring_and_controlling/server/handlers"
	"remote_monitoring_and_controlling/server/middleware"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type HTTPServer struct {
	router *echo.Echo

	handlers *handlers.Handler
	auth     *jwt.JWT

	appAddr      string
	permittedIPs []string
}

func NewHTTPServer(
	cfg *config.Config,
	handlers *handlers.Handler,
	auth *jwt.JWT,
) *HTTPServer {
	return &HTTPServer{
		appAddr:      cfg.AppAddr,
		permittedIPs: cfg.PermittedIPs,
		handlers:     handlers,
		auth:         auth,
	}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	s.router = echo.New()
	s.router.Use(middleware.LoggerMiddleware(), middleware.CorsMiddleware(), middleware.RecoverMiddleware())

	s.router.POST("/login", s.handlers.Login)
	s.router.POST("/2fa", s.handlers.TwoFA)

	// BAS (1C)
	bas := s.router.Group("/internal", middleware.CheckIPMiddleware(s.permittedIPs))
	bas.POST("/register", s.handlers.Register)
	bas.POST("/add_phone", s.handlers.AddAlternativeNumber)
	bas.POST("/tasks", s.handlers.SaveTasks)
	bas.POST("/vacations", s.handlers.SaveVacations)
	bas.POST("/payments", s.handlers.SavePayments)

	// User
	private := s.router.Group("", middleware.AuthMiddleware(s.auth))
	private.GET("/users", s.handlers.GetUsers)
	private.GET("/tasks", s.handlers.GetTasks)
	private.GET("/vacations", s.handlers.GetVacations)
	private.GET("/payments", s.handlers.GetPayments)
	private.POST("/logout", s.handlers.Logout)
	private.PUT("/task/status", s.handlers.UpdateTaskStatus)
	private.PATCH("/enable_2fa", s.handlers.EnableTwoFA)
	private.PATCH("/disable_2fa", s.handlers.DisableTwoFA)

	go func() {
		if err := s.router.Start(s.appAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Str("function", "Start").Err(err).Msg("Server start error")
		}
	}()

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
