package probes

import (
	"go.uber.org/fx"
	"user-service/internal/probes/delivery"
)

var Module = fx.Module(
	"Probes",
	fx.Invoke(
		delivery.RegisterHandlers,
	),
)
