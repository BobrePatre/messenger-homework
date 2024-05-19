package infrastructure

import (
	"auth-service/internal/infrastructure/grpc"
	"auth-service/internal/infrastructure/http"
	"auth-service/internal/infrastructure/logging"
	"go.uber.org/fx"
	"log/slog"
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
		http.RunHttpServer,
		grpc.RunGrpcServer,
		func(logger *slog.Logger) {
			logger.Info("Infrastructure Initialized")
		},
	),
)
