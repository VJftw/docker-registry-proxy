package plugin

import (
	"context"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type AuthProviderGRPCPlugin struct {
	plugin.Plugin
	Impl dockerregistryproxyv1.AuthenticationProviderServer
}

func (p *AuthProviderGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	dockerregistryproxyv1.RegisterAuthenticationProviderServer(s, p.Impl)
	return nil
}

func (p *AuthProviderGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return dockerregistryproxyv1.NewAuthenticationProviderClient(c), nil
}

func ServeAuthProviderPlugin(impl interface{}) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			PluginTypeAuthProvider:  &AuthProviderGRPCPlugin{Impl: impl.(dockerregistryproxyv1.AuthenticationProviderServer)},
			PluginTypeConfiguration: &ConfigurationGRPCPlugin{Impl: impl.(dockerregistryproxyv1.ConfigurationServer)},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
