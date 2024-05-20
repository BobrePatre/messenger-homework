package infrastructure

import (
	"go.uber.org/fx"
	"log/slog"
	"messaging-service/internal/infrastructure/grpc"
	"messaging-service/internal/infrastructure/http"
	"messaging-service/internal/infrastructure/logging"
	"messaging-service/internal/infrastructure/redis"
	"messaging-service/internal/infrastructure/validate"
	webAuthProvider "messaging-service/internal/infrastructure/web_auth_provider"
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
