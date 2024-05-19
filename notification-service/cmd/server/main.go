package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"notification-service/internal/infrastructure/logging"
)

func main() {
	app := fx.New(
		// Logger Deps
		fx.Provide(
			logging.LoggerConfig,
			logging.Logger,
		),

		// Configure logger for uber fx
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{
				Logger: logger,
			}
		}),

		// Infrastructure
		fx.Provide(),

		// Register handlers and other ....
		fx.Invoke(),

		// EntryPoint
		fx.Invoke(),
	)

	app.Run()
}
