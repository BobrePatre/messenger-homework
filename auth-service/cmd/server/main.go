package main

import (
	"auth-service/internal/infrastructure"
	"auth-service/internal/user"
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	app := fx.New(

		infrastructure.Module,
		user.Module,

		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{
				Logger: logger,
			}
		}),
	)

	app.Run()
}
