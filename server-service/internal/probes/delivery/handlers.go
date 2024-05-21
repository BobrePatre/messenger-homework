package delivery

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
	infraGrpc "server-service/internal/infrastructure/grpc"
	"server-service/internal/infrastructure/web_auth_provider/provider"
	"server-service/pkg/api/grpc/golang/probes"
	"strings"
)

type Handlers struct {
	logger       *slog.Logger
	authProvider provider.WebAuthProvider
	probes.UnimplementedProbeServiceServer
}

func RegisterHandlers(srv *grpc.Server, srvCfg *infraGrpc.Config, gw *runtime.ServeMux, logger *slog.Logger, authProvider provider.WebAuthProvider) error {
	ctx := context.Background()
	impl := &Handlers{logger: logger, authProvider: authProvider}

	probes.RegisterProbeServiceServer(srv, impl)
	err := probes.RegisterProbeServiceHandlerFromEndpoint(ctx, gw, srvCfg.Address(), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		return err
	}
	return nil
}

func (h *Handlers) Healthz(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (h *Handlers) Readyz(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	h.logger.Warn("MD", "md", md["authorization"])
	user, err := h.authProvider.Authorize(ctx, strings.Split(md["authorization"][0], " ")[1], nil)
	h.logger.Info("User authorized", "user", user)
	if err != nil {
		return nil, err
	}
	return nil, status.Error(codes.Unavailable, "service not ready")
}
