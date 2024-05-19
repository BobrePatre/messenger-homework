package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	grpcImpl "messaging-service/internal/adapters/primary/grpc"
	"messaging-service/internal/adapters/secondary/in_memory"
	"messaging-service/internal/application/interactors"
	grpcInfra "messaging-service/internal/infrastructure/grpc"
	"messaging-service/internal/infrastructure/http"
	"messaging-service/internal/infrastructure/logging"
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
			grpcImpl.RegisterMessagingHandlers,
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
