package user

import (
	"auth-service/internal/user/delivery"
	"auth-service/internal/user/repository"
	"auth-service/internal/user/service/interactors"
	"log/slog"

	"go.uber.org/fx"
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
