package user

import (
	"auth-service/internal/user/delivery"
	"auth-service/internal/user/repository"
	"auth-service/internal/user/service/interactors"
	"go.uber.org/fx"
	"log/slog"
)

var Module = fx.Module(
	"User Domain",
	fx.Provide(
		fx.Annotate(
			repository.NewUserRepository,
			fx.As(new(interactors.UserRepository)),
		),
		interactors.NewUserInteractor,
	),

	fx.Invoke(
		delivery.RegisterUserHandlers,
		func(logger *slog.Logger) {
			logger.Info("User Domain Connected")
		},
	),
)
