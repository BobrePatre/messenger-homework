package provider

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"time"
)

type Config struct {
	ClientId   string     `json:"client_id" env:"CLIENT_ID" validate:"required"`
	JwkOptions JwkOptions `json:"jwk_options" env-prefix:"JWK_OPTIONS_" validate:"required"`
}

type JwkOptions struct {
	RefreshJwkTimeout time.Duration `json:"refresh_timeout" env:"REFRESH_TIMEOUT" validate:"required"`
	JwkPublicUri      string        `json:"public_uri" env:"URI" validate:"required"`
}

func LoadConfig(validator *validator.Validate, logger *slog.Logger) (*Config, error) {
	var cfg struct {
		Config Config `json:"auth" env-prefix:"AUTH_"`
	}
	err := cleanenv.ReadConfig("config.json", &cfg)
	if err != nil {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
	}
	err = validator.Struct(cfg.Config)
	if err != nil {
		logger.Error("Failed to validate security config", "error", err)
		return nil, err
	}
	return &cfg.Config, nil
}
