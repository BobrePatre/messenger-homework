package grpc

import (
	"google.golang.org/grpc"
	"log/slog"
	"server-service/pkg/api/grpc/golang/server"
)

type AuthHandler struct {
	logger *slog.Logger
	server.UnimplementedServerServiceServer
}

func RegisterServerHandlers(srv *grpc.Server, logger *slog.Logger) error {

	impl := &AuthHandler{logger: logger}

	server.RegisterServerServiceServer(srv, impl)

	return nil
}
