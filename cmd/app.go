package cmd

import (
	"remote_monitoring_and_controlling/config"
	"remote_monitoring_and_controlling/internal/user"
	"remote_monitoring_and_controlling/pkg/bcrypt"
	"remote_monitoring_and_controlling/pkg/jwt"
	"remote_monitoring_and_controlling/pkg/otp"
	"remote_monitoring_and_controlling/pkg/postgres"
	"remote_monitoring_and_controlling/server"
	"remote_monitoring_and_controlling/server/handlers"

	"go.uber.org/fx"
)

func Exec() fx.Option {
	return fx.Options(
		fx.Provide(
			config.Get,
			bcrypt.NewBcrypt,
			otp.NewOTP,

			postgres.NewPostgres,
			fx.Annotate(
				annotationDupl[postgres.Postgres],
				fx.As(new(user.DB)),
			),

			jwt.NewJWT,
			user.NewUserService,
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

func prepareHooks(server *server.HTTPServer, postgres *postgres.Postgres, lc fx.Lifecycle) {
	lc.Append(fx.Hook{OnStart: postgres.Start, OnStop: postgres.Stop})
	lc.Append(fx.Hook{OnStart: server.Start, OnStop: server.Stop})
}
