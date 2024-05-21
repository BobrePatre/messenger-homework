package keycloak_redis

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"messaging-service/internal/infrastructure/web_auth_provider/provider"
	models2 "messaging-service/internal/infrastructure/web_auth_provider/provider/models"
)

var _ provider.WebAuthProvider = (*Provider)(nil)

type Provider struct {
	redis    *redis.Client
	jwkOpts  provider.JwkOptions
	validate *validator.Validate
	logger   *slog.Logger
	clientID string
}

func NewProvider(redis *redis.Client, validate *validator.Validate, cfg *provider.Config, logger *slog.Logger) *Provider {
	return &Provider{
		redis:    redis,
		jwkOpts:  cfg.JwkOptions,
		validate: validate,
		clientID: cfg.ClientId,
		logger:   logger,
	}
}

func (p *Provider) Authorize(ctx context.Context, tokenString string, neededRoles []string) (models2.UserDetails, error) {
	token, err := p.VerifyToken(ctx, tokenString)
	if err != nil {
		p.logger.Error("failed to verify token", slog.String("err", err.Error()))
		return models2.UserDetails{}, models2.InvalidTokenError
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !(ok && token.Valid) {
		p.logger.Error("failed to get claims")
		return models2.UserDetails{}, models2.InvalidTokenError
	}

	if claims["sub"] == "" || claims["sub"] == nil {
		p.logger.Error("failed to validate sub claim")
		return models2.UserDetails{}, models2.InvalidTokenError
	}

	err = p.validate.Var(claims["sub"], "uuid4")
	if err != nil {
		p.logger.Error("failed to validate sub claim", slog.String("err", err.Error()))
		return models2.UserDetails{}, err
	}

	var userRoles []string
	if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
		if authClient, ok := resourceAccess[p.clientID].(map[string]interface{}); ok {
			if err := mapstructure.Decode(authClient["roles"], &userRoles); err != nil {
				p.logger.Error("cannot get user roles", slog.String("err", err.Error()))
				userRoles = []string{}
			}
		}
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		userEmail = ""
	}

	userDetails := models2.UserDetails{
		Roles:      userRoles,
		UserId:     claims["sub"].(string),
		Email:      userEmail,
		Username:   claims["preferred_username"].(string),
		Name:       claims["name"].(string),
		FamilyName: claims["family_name"].(string),
	}

	if !p.IsUserHaveRoles(neededRoles, userRoles) {
		p.logger.Info("user data", slog.Any("userDetails", userDetails))
		p.logger.Error("user doesn't have needed roles", slog.Any("neededRoles", neededRoles), slog.Any("userRoles", userRoles))
		return userDetails, models2.AccessDeniedError
	}

	return userDetails, nil
}
