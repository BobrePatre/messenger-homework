package infrastructure

import (
	"go.uber.org/fx"
	"log/slog"
	"messaging-service/internal/infrastructure/grpc"
	"messaging-service/internal/infrastructure/http"
	"messaging-service/internal/infrastructure/logging"
)

var Module = fx.Module(
	"Infrastructure",

	// Logger
	fx.Provide(
		logging.LoadConfig,
		logging.Logger,
	),

	// Grpc
	fx.Provide(
		grpc.LoadConfig,
		grpc.NewGrpcServer,
	),

	// Http
	fx.Provide(
		http.LoadConfig,
		http.NewGatewayServer,
		http.NewHttpServer,
	),

	// Module Entrypoint
	fx.Invoke(
		grpc.RunGrpcServer,
		http.RunHttpServer,
		func(logger *slog.Logger) {
			logger.Info("Infrastructure Initialized")
		},
	),
)
