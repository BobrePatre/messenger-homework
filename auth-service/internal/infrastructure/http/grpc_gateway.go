package http

import (
	"auth-service/pkg/api/grpc/golang/auth"
	"auth-service/pkg/api/grpc/golang/probes"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGatewayServer() (*runtime.ServeMux, error) {
	gatewayMux := runtime.NewServeMux()
	ctx := context.Background()

	err := auth.RegisterAuthServiceHandlerFromEndpoint(ctx, gatewayMux, "0.0.0.0:2000", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		return nil, err
	}

	err = probes.RegisterProbeServiceHandlerFromEndpoint(ctx, gatewayMux, "0.0.0.0:2000", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		return nil, err
	}

	return gatewayMux, nil
}
