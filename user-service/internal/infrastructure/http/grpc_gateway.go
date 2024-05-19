package http

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log/slog"
	"user-service/pkg/api/grpc/golang/probes"
	"user-service/pkg/api/grpc/golang/users"
)

func NewGatewayServer(logger *slog.Logger) (*runtime.ServeMux, error) {
	gatewayMux := runtime.NewServeMux()
	ctx := context.Background()

	err := users.RegisterUsersServiceHandlerFromEndpoint(ctx, gatewayMux, "0.0.0.0:2000", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = probes.RegisterProbeServiceHandlerFromEndpoint(ctx, gatewayMux, "0.0.0.0:2000", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	return gatewayMux, nil
}
