package grpc

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"time"
)

var grpcServerTag = slog.String("server", "grpc_server")

type Config struct {
	Host string `json:"host" env-default:"0.0.0.0"`
	Port int    `json:"port" env-default:"2000"`
}

func LoadConfig() (*Config, error) {
	var cfg struct {
		Config Config `json:"grpc"`
	}
	err := cleanenv.ReadConfig("config.json", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg.Config, nil
}

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

func RunGrpcServer(lc fx.Lifecycle, srv *grpc.Server, logger *slog.Logger, cfg *Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting server", grpcServerTag, "address", net.JoinHostPort(cfg.Host, fmt.Sprint(cfg.Port)))
			listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Host, fmt.Sprint(cfg.Port)))
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
