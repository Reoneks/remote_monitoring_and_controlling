package cmd

import (
	"project/config"
	"project/internal/user"
	"project/pkg/bcrypt"
	"project/pkg/jwt"
	"project/pkg/otp"
	"project/pkg/postgres"
	"project/server"
	"project/server/handlers"

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
