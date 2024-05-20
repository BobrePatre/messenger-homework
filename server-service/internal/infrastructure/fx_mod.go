package infrastructure

import (
	"go.uber.org/fx"
	"log/slog"
	"server-service/internal/infrastructure/grpc"
	"server-service/internal/infrastructure/http"
	"server-service/internal/infrastructure/logging"
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
