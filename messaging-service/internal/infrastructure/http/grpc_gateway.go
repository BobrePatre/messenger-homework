package http

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log/slog"
	"messaging-service/pkg/api/grpc/golang/messaging"
	"messaging-service/pkg/api/grpc/golang/probes"
)

func NewGatewayServer(logger *slog.Logger) (*runtime.ServeMux, error) {
	gatewayMux := runtime.NewServeMux()
	ctx := context.Background()

	err := messaging.RegisterMessagingHandlerFromEndpoint(ctx, gatewayMux, "0.0.0.0:2000", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	err = probes.RegisterProbeServiceHandlerFromEndpoint(ctx, gatewayMux, "0.0.0.0:2000", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return nil, err
	}

	return gatewayMux, nil
}
