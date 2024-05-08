package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-faker/faker/v4"
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
	"user-service/pkg/api/grpc/golang/probes"
	"user-service/pkg/api/grpc/golang/users"
)

type Impl struct {
	probes.UnimplementedProbeServiceServer
	users.UnimplementedUsersServiceServer
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
	err = users.RegisterUsersServiceHandlerFromEndpoint(ctx, gatewayServer, "localhost:2000", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	httpServer.Any("/*any", echo.WrapHandler(gatewayServer))

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	reflection.Register(grpcServer)

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

func (impl *Impl) Healthz(context.Context, *emptypb.Empty) (*probes.ProbeResponse, error) {
	return &probes.ProbeResponse{
		Message: "ok, im user service",
		Status:  probes.Status_OK,
	}, nil
}

func (impl *Impl) Readyz(context.Context, *emptypb.Empty) (*probes.ProbeResponse, error) {
	return &probes.ProbeResponse{
		Message: "ok, im user service and im not ready",
		Status:  probes.Status_UNAVAILABLE,
	}, nil
}

func (impl *Impl) AddFriend(context.Context, *users.AddFriendRequest) (*users.AddFriendResponse, error) {
	return &users.AddFriendResponse{
		RequestId: uuid.NewString(),
	}, nil
}
func (impl *Impl) SearchUser(context.Context, *users.SearchUserRequst) (*users.SearchUserResponse, error) {
	randomUsername := faker.Username()
	randomAvatarId := uuid.NewString()

	return &users.SearchUserResponse{
		Users: []*users.UserData{
			{
				UserId:   uuid.New().String(),
				Username: &randomUsername,
				AvatarId: &randomAvatarId,
			},
		},
	}, nil
}
func (impl *Impl) RemoveFriend(context.Context, *users.RemoveFriendRequest) (*users.RemoveFriendResponse, error) {
	return &users.RemoveFriendResponse{
		Status: users.Status_SUCCESS,
	}, nil
}
func (impl *Impl) ProcessFriendRequest(context.Context, *users.ProcessFriendRequestRequest) (*users.ProccessFriendRequestResponse, error) {
	return &users.ProccessFriendRequestResponse{
		Status: users.Status_ERROR,
	}, nil
}
func (impl *Impl) GetFriends(context.Context, *emptypb.Empty) (*users.Friends, error) {
	return &users.Friends{
		Friends: []*users.Friend{
			{
				UserId:       uuid.NewString(),
				FriendStatus: users.FriendStatus_PENDING,
			},
		},
	}, nil
}
