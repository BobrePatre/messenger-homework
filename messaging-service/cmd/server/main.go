package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"messaging-service/internal/infrastructure"
)

func main() {
	app := fx.New(

		infrastructure.Module,

		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{
				Logger: logger,
			}
		}),
	)

	app.Run()
}
