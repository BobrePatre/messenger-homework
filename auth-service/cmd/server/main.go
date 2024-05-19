package main

import (
	grpcImpl "auth-service/internal/adapters/primary/grpc"
	"auth-service/internal/adapters/secondary/in_memory"
	"auth-service/internal/application/interactors"
	grpcInfra "auth-service/internal/infrastructure/grpc"
	"auth-service/internal/infrastructure/http"
	"auth-service/internal/infrastructure/logging"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
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
		fx.Provide(
			grpcInfra.NewGrpcServer,
			http.NewHttpServer,
			http.NewGatewayServer,
		),

		// User
		fx.Provide(
			fx.Annotate(
				in_memory.NewUserRepository,
				fx.As(new(interactors.UserRepository)),
			),
			interactors.NewUserInteractor,
		),

		// Register handlers and other ....
		fx.Invoke(
			grpcImpl.RegisterAuthHandlers,
			http.NewGatewayServer,
		),

		// EntryPoint
		fx.Invoke(
			http.RunHttpServer,
			grpcInfra.RunGrpcServer,
		),
	)

	app.Run()
}
