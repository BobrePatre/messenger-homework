package grpc

import (
	"auth-service/pkg/api/grpc/golang/auth"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type AuthHandler struct {
	logger *slog.Logger
	auth.UnimplementedAuthServiceServer
}

func RegisterAuthHandlers(srv *grpc.Server, logger *slog.Logger) error {

	impl := &AuthHandler{logger: logger}

	auth.RegisterAuthServiceServer(srv, impl)

	return nil
}

func (h *AuthHandler) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	h.logger.Info("New Login Req",
		"Login", request.Login,
		"Password", request.Password,
	)
	return &auth.LoginResponse{
		Token: uuid.NewString(),
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	h.logger.Info("New Reg Req",
		"Username", request.Username,
		"Email", request.Email,
		"Password", request.Password,
		"PasswordConfirmation", request.PasswordConfirmation,
	)
	if request.Password != request.PasswordConfirmation {
		h.logger.Warn("Password mismatch",
			"Password", request.Password,
			"PasswordConfirmation", request.PasswordConfirmation,
		)
		return nil, status.Error(codes.InvalidArgument, "password mismatch")
	}
	return &auth.RegisterResponse{
		Token: uuid.NewString(),
	}, nil
}
