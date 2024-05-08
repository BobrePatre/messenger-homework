package main

import (
	"context"
	"errors"
	"fmt"
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
	"server-service/pkg/api/grpc/golang/probes"
	"server-service/pkg/api/grpc/golang/server"
	"sync"
	"syscall"
)

type Impl struct {
	probes.UnimplementedProbeServiceServer
	server.UnimplementedServerServiceServer
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
	err = server.RegisterServerServiceHandlerFromEndpoint(ctx, gatewayServer, "localhost:2000", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	httpServer.Any("/*any", echo.WrapHandler(gatewayServer))

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	reflection.Register(grpcServer)

	probes.RegisterProbeServiceServer(grpcServer, impl)
	server.RegisterServerServiceServer(grpcServer, impl)
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
		Message: "ok, im server service",
		Status:  probes.Status_OK,
	}, nil
}

func (impl *Impl) Readyz(context.Context, *emptypb.Empty) (*probes.ProbeResponse, error) {
	return &probes.ProbeResponse{
		Message: "ok, im server service and im not ready",
		Status:  probes.Status_UNAVAILABLE,
	}, nil
}

func (impl *Impl) CreateServer(ctx context.Context, request *server.CreateServerRequest) (*server.CreateServerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (impl *Impl) SearchServer(ctx context.Context, request *server.SearchServerRequest) (*server.SearchServerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (impl *Impl) SubsribeServer(ctx context.Context, request *server.SubscribeServerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (impl *Impl) UnSubsribeServer(ctx context.Context, request *server.UnSubscribeServerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (impl *Impl) GetServers(ctx context.Context, empty *emptypb.Empty) (*server.GetServersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (impl *Impl) CreateInvite(ctx context.Context, request *server.CreateServerRequest) (*server.CreateServerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (impl *Impl) AcceptInvite(ctx context.Context, request *server.AcceptInviteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (impl *Impl) mustEmbedUnimplementedServerServiceServer() {
	//TODO implement me
	panic("implement me")
}
