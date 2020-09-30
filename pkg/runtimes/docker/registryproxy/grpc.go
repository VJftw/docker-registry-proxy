package registryproxy

import (
	"context"
	"fmt"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/google/go-containerregistry/pkg/authn"
	"go.uber.org/zap"
)

// RegistryAuthResolveGRPC resolves the registry authentication from via GRPC
func RegistryAuthResolveGRPC(client dockerregistryproxyv1.AuthenticationProviderClient) (authn.Authenticator, error) {
	resp, err := client.Provide(context.Background(), &dockerregistryproxyv1.ProvideRequest{})
	if err != nil {
		return nil, err
	}

	cfg := authn.AuthConfig{
		Username:      resp.GetUsername(),
		Password:      resp.GetPassword(),
		Auth:          resp.GetAuth(),
		IdentityToken: resp.GetIdentityToken(),
		RegistryToken: resp.GetRegistryToken(),
	}

	logger.Info("recieved credentials from auth provider", zap.String("credentials", fmt.Sprintf("%#v", cfg)))

	return authn.FromConfig(cfg), nil
}
