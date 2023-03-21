package cmd

import (
	"project/config"
	"project/pkg/bcrypt"
	"project/pkg/jwt"
	"project/pkg/postgres"
	"project/pkg/telegram"
	"project/server"
	"project/server/handlers"
	"project/server/middleware"

	"go.uber.org/fx"
)

func Exec() fx.Option {
	return fx.Options(
		fx.Provide(
			config.Get,

			bcrypt.NewBcrypt,
			fx.Annotate(
				annotationDupl[postgres.Postgres],
				fx.As(new(jwt.Bcrypt)),
				fx.As(new(handlers.Bcrypt)),
			),

			postgres.NewPostgres,
			fx.Annotate(
				annotationDupl[postgres.Postgres],
				fx.As(new(handlers.Postgres)),
				fx.As(new(telegram.Postgres)),
			),

			telegram.NewTelegram,
			fx.Annotate(
				annotationDupl[telegram.Telegram],
				fx.As(new(handlers.Telegram)),
			),

			fx.Annotate(
				jwt.NewJWT,
				fx.As(new(middleware.Auth)),
			),

			handlers.NewHandler,
			server.NewHTTPServer,
		),
		fx.Invoke(
			prepareHooks,
		),
	)
}

func annotationDupl[T any](v *T) *T {
	return v
}

func prepareHooks(server *server.HTTPServer, postgres *postgres.Postgres, telegram *telegram.Telegram, lc fx.Lifecycle) {
	lc.Append(fx.Hook{OnStart: postgres.Start, OnStop: postgres.Stop})
	lc.Append(fx.Hook{OnStart: telegram.Start, OnStop: telegram.Stop})
	lc.Append(fx.Hook{OnStart: server.Start, OnStop: server.Stop})
}
