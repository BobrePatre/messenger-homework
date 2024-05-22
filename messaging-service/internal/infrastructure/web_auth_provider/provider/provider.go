package provider

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"google.golang.org/grpc"
	"messaging-service/internal/infrastructure/web_auth_provider/provider/models"
)

type WebAuthProvider interface {
	VerifyToken(ctx context.Context, tokenString string) (*jwt.Token, error)
	TokenKeyfunc(ctx context.Context) jwt.Keyfunc
	FetchJwkSet(ctx context.Context) (jwk.Set, error)
	IsUserHaveRoles(roles []string, userRoles []string) bool
	SerializeJwkSet(key jwk.Set) (string, error)
	DeserializeJwkSet(serializedKey string) (jwk.Set, error)
	Authorize(ctx context.Context, tokenString string, roles []string) (models.UserDetails, error)
}

const (
	JwkKeySet      = "jwk-set"
	UserDetailsKey = "user-details"
)

type (
	AuthHttpMiddleware              func(roles ...string) echo.HandlerFunc
	AuthHttpMiddlewareWrapper       func(provider WebAuthProvider) func(roles ...string) echo.HandlerFunc
	AuthGrpcUnaryInterceptorWrapper func(provider WebAuthProvider) grpc.UnaryServerInterceptor
)
