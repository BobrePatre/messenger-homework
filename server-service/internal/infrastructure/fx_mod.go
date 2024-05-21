package infrastructure

import (
	"go.uber.org/fx"
	"log/slog"
	"server-service/internal/infrastructure/grpc"
	"server-service/internal/infrastructure/http"
	"server-service/internal/infrastructure/logging"
	"server-service/internal/infrastructure/redis"
	"server-service/internal/infrastructure/validate"
	webAuthProvider "server-service/internal/infrastructure/web_auth_provider"
)

var Module = fx.Module(
	"Infrastructure",

	// Validator
	fx.Provide(
		validate.NewValidate,
	),

	// Logger
	fx.Provide(
		logging.LoadConfig,
		logging.Logger,
	),

	// Redis
	fx.Provide(
		redis.LoadConfig,
		redis.NewClient,
	),

	// Security
	webAuthProvider.Module,

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
