package webAuthProvider

import (
	"context"
	"go.uber.org/fx"
	"log/slog"
	"server-service/internal/infrastructure/web_auth_provider/keycloak_redis"
	"server-service/internal/infrastructure/web_auth_provider/provider"
)

var Module = fx.Module(
	"Web Secutiry",
	fx.Provide(
		provider.LoadConfig,
		fx.Annotate(
			keycloak_redis.NewProvider,
			fx.As(new(provider.WebAuthProvider)),
		),
	),
	fx.Invoke(
		func(authProvider provider.WebAuthProvider, logger *slog.Logger) error {
			_, err := authProvider.FetchJwkSet(context.Background())
			if err != nil {
				return err
			}
			logger.Info("jwks loaded")
			return nil
		},
	),
)
