package grpc

import (
	"google.golang.org/grpc"
	"log/slog"
	"messaging-service/pkg/api/grpc/golang/messaging"
)

type MessagingHandlers struct {
	logger *slog.Logger
	messaging.UnimplementedMessagingServer
}

func RegisterMessagingHandlers(srv *grpc.Server, logger *slog.Logger) error {

	impl := &MessagingHandlers{logger: logger}

	messaging.RegisterMessagingServer(srv, impl)

	return nil
}
