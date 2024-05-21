package probes

import (
	"go.uber.org/fx"
	"server-service/internal/probes/delivery"
)

var Module = fx.Module(
	"Probes",
	fx.Invoke(
		delivery.RegisterHandlers,
	),
)
