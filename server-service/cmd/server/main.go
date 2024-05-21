package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"server-service/internal/infrastructure"
	"server-service/internal/probes"
)

func main() {
	app := fx.New(

		infrastructure.Module,
		probes.Module,

		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{
				Logger: logger,
			}
		}),
	)

	app.Run()
}
