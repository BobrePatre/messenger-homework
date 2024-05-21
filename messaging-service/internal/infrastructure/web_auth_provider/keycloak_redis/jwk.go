package keycloak_redis

import (
	"context"
	"github.com/lestrrat-go/jwx/jwk"
	webAuthProvider "messaging-service/internal/infrastructure/web_auth_provider/provider"
)

var _ webAuthProvider.WebAuthProvider = (*Provider)(nil)

func (p *Provider) FetchJwkSet(ctx context.Context) (jwk.Set, error) {

	result, err := p.redis.Get(ctx, webAuthProvider.JwkKeySet).Result()
	if err == nil {
		p.logger.Info("Jwk get from cache")
		resultSet, err := p.DeserializeJwkSet(result)
		if err != nil {
			return nil, err
		}
		return resultSet, nil
	}

	p.logger.Info("Fetch Jwk from remote")
	resultSet, err := jwk.Fetch(ctx, p.jwkOpts.JwkPublicUri)
	if err != nil {
		p.logger.Error("Failed to fetch jwk public uri from remote", "error", err)
		return nil, err
	}

	p.logger.Info("Serializing jwks")
	serializedKeySet, err := p.SerializeJwkSet(resultSet)
	if err != nil {
		p.logger.Error("Failed to serialize jwk set", "error", err)
		return nil, err
	}

	err = p.redis.Set(ctx, webAuthProvider.JwkKeySet, serializedKeySet, p.jwkOpts.RefreshJwkTimeout).Err()
	if err != nil {
		p.logger.Error("Failed to set jwk set", "error", err)
		return nil, err
	}

	return resultSet, nil

}
