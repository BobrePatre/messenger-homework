package grpc

import (
	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"time"
)

var grpcServerTag = slog.String("server", "grpc_server")

const (
	serverPort = "2000"
	serverHost = "0.0.0.0"
)

func NewGrpcServer(logger *slog.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				logger.Info("New Grpc Req",
					"Method", info.FullMethod,
					"Server", info.Server,
				)
				return handler(ctx, req)
			},
		),
		grpc.Creds(insecure.NewCredentials()),
	)
	reflection.Register(server)
	return server
}

func RunGrpcServer(lc fx.Lifecycle, srv *grpc.Server, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting server", grpcServerTag, "address", net.JoinHostPort(serverHost, serverPort))
			listener, err := net.Listen("tcp", net.JoinHostPort(serverHost, serverPort))
			if err != nil {
				logger.Error("cannot start server", "error", err.Error(), grpcServerTag)
				return err
			}
			go func() {
				err := srv.Serve(listener)
				if err != nil {
					logger.Error("cannot start server", "error", err.Error(), grpcServerTag)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down", grpcServerTag)
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			srv.GracefulStop()
			return nil
		},
	})
}
