package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	grpcImpl "server-service/internal/adapters/primary/grpc"
	"server-service/internal/adapters/secondary/in_memory"
	"server-service/internal/application/interactors"
	grpcInfra "server-service/internal/infrastructure/grpc"
	"server-service/internal/infrastructure/http"
	"server-service/internal/infrastructure/logging"
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
			grpcImpl.RegisterServerHandlers,
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
