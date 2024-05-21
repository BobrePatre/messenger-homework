package probes

import (
	"go.uber.org/fx"
	"messaging-service/internal/probes/delivery"
)

var Module = fx.Module(
	"Probes",
	fx.Invoke(
		delivery.RegisterHandlers,
	),
)
