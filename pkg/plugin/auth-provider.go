package plugin

import (
	"context"

	v1 "github.com/VJftw/docker-registry-proxy/pkg/genproto/v1"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type AuthProviderGRPCPlugin struct {
	plugin.Plugin
	Impl v1.AuthenticationProviderServer
}

func (p *AuthProviderGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	v1.RegisterAuthenticationProviderServer(s, p.Impl)
	return nil
}

func (p *AuthProviderGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return v1.NewAuthenticationProviderClient(c), nil
}

func ServeAuthProviderPlugin(impl interface{}) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			PluginTypeAuthProvider:  &AuthProviderGRPCPlugin{Impl: impl.(v1.AuthenticationProviderServer)},
			PluginTypeConfiguration: &ConfigurationGRPCPlugin{Impl: impl.(v1.ConfigurationServer)},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
