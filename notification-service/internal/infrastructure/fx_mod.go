package infrastructure

import (
	"go.uber.org/fx"
	"log/slog"
	"notification-service/internal/infrastructure/logging"
)

var Module = fx.Module(
	"Infrastructure",

	// Logger
	fx.Provide(
		logging.LoadConfig,
		logging.Logger,
	),

	// TODO: Add kafka and WS

	// Module Entrypoint
	fx.Invoke(
		// TODO: Add entrypoit for kafka and WS
		func(logger *slog.Logger) {
			logger.Info("Infrastructure Initialized")
		},
	),
)
