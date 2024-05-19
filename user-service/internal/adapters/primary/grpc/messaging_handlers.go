package grpc

import (
	"google.golang.org/grpc"
	"log/slog"
	"user-service/pkg/api/grpc/golang/users"
)

type MessagingHandlers struct {
	logger *slog.Logger
	users.UnimplementedUsersServiceServer
}

func RegisterMessagingHandlers(srv *grpc.Server, logger *slog.Logger) error {

	impl := &MessagingHandlers{logger: logger}

	users.RegisterUsersServiceServer(srv, impl)

	return nil
}
