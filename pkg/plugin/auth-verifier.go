package plugin

import (
	"context"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type AuthVerifierGRPCPlugin struct {
	plugin.Plugin
	Impl dockerregistryproxyv1.AuthenticationVerifierAPIServer
}

func (p *AuthVerifierGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	dockerregistryproxyv1.RegisterAuthenticationVerifierAPIServer(s, p.Impl)
	return nil
}

func (p *AuthVerifierGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return dockerregistryproxyv1.NewAuthenticationVerifierAPIClient(c), nil
}

func ServeAuthVerifierPlugin(impl interface{}) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			PluginTypeAuthVerifier:  &AuthVerifierGRPCPlugin{Impl: impl.(dockerregistryproxyv1.AuthenticationVerifierAPIServer)},
			PluginTypeConfiguration: &ConfigurationGRPCPlugin{Impl: impl.(dockerregistryproxyv1.ConfigurationAPIServer)},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
