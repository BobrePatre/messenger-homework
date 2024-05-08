package main

import (
	"auth-service/pkg/api/grpc/golang/auth"
	"auth-service/pkg/api/grpc/golang/probes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Impl struct {
	auth.UnimplementedAuthServiceServer
	probes.UnimplementedProbeServiceServer
}

func main() {

	ctx := context.Background()
	impl := NewImpl()
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	gatewayServer := runtime.NewServeMux()
	httpServer := echo.New()
	httpServer.Use(middleware.Recover())
	httpServer.Use(middleware.Logger())
	err := auth.RegisterAuthServiceHandlerFromEndpoint(ctx, gatewayServer, "localhost:2000", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		panic(err)
	}
	err = probes.RegisterProbeServiceHandlerFromEndpoint(ctx, gatewayServer, "localhost:2000", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		panic(err)
	}

	httpServer.Any("/*any", echo.WrapHandler(gatewayServer))

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	reflection.Register(grpcServer)

	auth.RegisterAuthServiceServer(grpcServer, impl)
	probes.RegisterProbeServiceServer(grpcServer, impl)

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()

		listener, err := net.Listen("tcp", ":2000")
		if err != nil {
			panic(err)
		}

		err = grpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for _, route := range httpServer.Routes() {
			fmt.Println(*route)
		}
		err := httpServer.Start(":8080")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-stopChan
	grpcServer.GracefulStop()
	err = httpServer.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
	wg.Wait()
	fmt.Println("done")
}

func NewImpl() *Impl {
	return &Impl{}
}

func (s *Impl) Login(_ context.Context, _ *auth.LoginRequest) (*auth.LoginResponse, error) {
	return &auth.LoginResponse{
		Token: uuid.NewString(),
	}, nil
}

func (s *Impl) Register(_ context.Context, _ *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return &auth.RegisterResponse{
		Token: uuid.NewString(),
	}, nil
}

func (s *Impl) Healthz(_ context.Context, _ *emptypb.Empty) (*probes.ProbeResponse, error) {
	return &probes.ProbeResponse{
		Status:  probes.Status_OK,
		Message: "ok, im auth service",
	}, nil
}

func (s *Impl) Readyz(_ context.Context, _ *emptypb.Empty) (*probes.ProbeResponse, error) {
	return &probes.ProbeResponse{
		Status:  probes.Status_UNAVAILABLE,
		Message: "im auth service and im not ready",
	}, nil
}
