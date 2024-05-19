package delivery

import (
	"auth-service/pkg/api/grpc/golang/auth"
	"google.golang.org/grpc"
	"log/slog"
)

type UserHandlers struct {
	logger *slog.Logger
	auth.UnimplementedAuthServiceServer
}

func RegisterUserHandlers(srv *grpc.Server, logger *slog.Logger) error {

	impl := &UserHandlers{logger: logger}

	auth.RegisterAuthServiceServer(srv, impl)

	return nil
}
