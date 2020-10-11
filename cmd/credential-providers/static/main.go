package main

import (
	"context"
	"log"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
)

const (
	flagUsername = "username"
	flagPassword = "password"
)

func main() {
	plugin.ServeAuthProviderPlugin(&Provider{})
}

// Provider represents an AuthenticationProvider using static credentials
type Provider struct {
	dockerregistryproxyv1.AuthenticationProviderAPIServer
	dockerregistryproxyv1.ConfigurationAPIServer

	username string
	password string
}

// Provide returns credentials
func (p *Provider) Provide(ctx context.Context, req *dockerregistryproxyv1.ProvideRequest) (*dockerregistryproxyv1.ProvideResponse, error) {
	return &dockerregistryproxyv1.ProvideResponse{
		Username: p.username,
		Password: p.password,
	}, nil
}

// GetConfigurationSchema returns the schema for the plugin
func (p *Provider) GetConfigurationSchema(ctx context.Context, _ *empty.Empty) (*dockerregistryproxyv1.GetConfigurationSchemaResponse, error) {
	return &dockerregistryproxyv1.GetConfigurationSchemaResponse{
		Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttribute{
			flagUsername: {
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
				Description:   "the static username",
			},
			flagPassword: {
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
				Description:   "the static password",
			},
		},
	}, nil
}

// Configure configures the plugin
func (p *Provider) Configure(ctx context.Context, req *dockerregistryproxyv1.ConfigureRequest) (*empty.Empty, error) {
	if val, ok := req.Attributes[flagUsername]; ok {
		username, err := plugin.UnmarshalConfigurationValue(dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING, val.GetValue())
		if err != nil {
			return nil, err
		}
		p.username = username.(string)
		log.Printf("configured username as '%s'", p.username)
	}
	if val, ok := req.Attributes[flagPassword]; ok {
		password, err := plugin.UnmarshalConfigurationValue(dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING, val.GetValue())
		if err != nil {
			return nil, err
		}
		p.password = password.(string)
		log.Printf("configured password")
	}
	return &empty.Empty{}, nil
}
