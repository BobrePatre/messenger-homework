package delivery

import (
	infraGrpc "auth-service/internal/infrastructure/grpc"
	"auth-service/pkg/api/grpc/golang/auth"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type UserHandlers struct {
	logger *slog.Logger
	auth.UnimplementedAuthServiceServer
}

func RegisterUserHandlers(srv *grpc.Server, gw *runtime.ServeMux, srvCfg *infraGrpc.Config, logger *slog.Logger) error {

	impl := &UserHandlers{logger: logger}
	ctx := context.Background()

	auth.RegisterAuthServiceServer(srv, impl)
	err := auth.RegisterAuthServiceHandlerFromEndpoint(ctx, gw, srvCfg.Address(), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		return err
	}

	return nil
}
