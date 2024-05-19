package http

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log/slog"
	"server-service/pkg/api/grpc/golang/server"
)

func NewGatewayServer(logger *slog.Logger) (*runtime.ServeMux, error) {
	gatewayMux := runtime.NewServeMux()
	ctx := context.Background()

	err := server.RegisterServerServiceHandlerFromEndpoint(ctx, gatewayMux, "0.0.0.0:2000", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	return gatewayMux, nil
}
