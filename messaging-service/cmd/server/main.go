package main

import (
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
	"messaging-service/pkg/api/grpc/golang/messaging"
	"messaging-service/pkg/api/grpc/golang/probes"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Impl struct {
	probes.UnimplementedProbeServiceServer
	messaging.UnimplementedMessagingServer
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
	err := probes.RegisterProbeServiceHandlerFromEndpoint(ctx, gatewayServer, "localhost:2000", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		panic(err)
	}
	err = messaging.RegisterMessagingHandlerFromEndpoint(ctx, gatewayServer, "localhost:2000", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	httpServer.Any("/*any", echo.WrapHandler(gatewayServer))

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	reflection.Register(grpcServer)

	probes.RegisterProbeServiceServer(grpcServer, impl)
	messaging.RegisterMessagingServer(grpcServer, impl)
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

func (impl *Impl) Healthz(context.Context, *emptypb.Empty) (*probes.ProbeResponse, error) {
	return &probes.ProbeResponse{
		Message: "ok, im messaging service",
		Status:  probes.Status_OK,
	}, nil
}

func (impl *Impl) Readyz(context.Context, *emptypb.Empty) (*probes.ProbeResponse, error) {
	return &probes.ProbeResponse{
		Message: "ok, im messaging service and im not ready",
		Status:  probes.Status_UNAVAILABLE,
	}, nil
}

func (impl *Impl) SendMessage(context.Context, *messaging.SendMessageRequest) (*messaging.SendMessageResponse, error) {
	return &messaging.SendMessageResponse{
		MessageId: uuid.NewString(),
	}, nil
}
