package provider

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	ClientId   string     `json:"client_id" env:"CLIENT_ID" validate:"required"`
	JwkOptions JwkOptions `json:"jwk_options" env:"JWK_OPTIONS" validate:"reuired"`
}

type JwkOptions struct {
	RefreshJwkTimeout time.Duration `json:"refresh_jwk_timeout" env:"JWK_TIMEOUT" validate:"required"`
	JwkPublicUri      string        `json:"jwk_public_uri" env:"JWK_URI" validate:"required"`
}

func LoadConfig() (*Config, error) {
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
	return &cfg.Config, nil
}
