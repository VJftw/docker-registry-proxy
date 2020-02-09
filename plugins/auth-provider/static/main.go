package main

import (
	"context"
	"log"

	v1 "github.com/VJftw/docker-registry-proxy/pkg/genproto/v1"
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
	v1.AuthenticationProviderServer
	v1.ConfigurationServer

	username string
	password string
}

// Provide returns credentials
func (p *Provider) Provide(ctx context.Context, req *v1.ProvideRequest) (*v1.ProvideResponse, error) {
	return &v1.ProvideResponse{
		Username: p.username,
		Password: p.password,
	}, nil
}

// GetConfigurationSchema returns the schema for the plugin
func (p *Provider) GetConfigurationSchema(ctx context.Context, _ *empty.Empty) (*v1.ConfigurationSchema, error) {
	return &v1.ConfigurationSchema{
		Attributes: map[string]*v1.ConfigurationAttribute{
			flagUsername: &v1.ConfigurationAttribute{
				AttributeType: v1.ConfigType_STRING,
				Description:   "the static username",
			},
			flagPassword: &v1.ConfigurationAttribute{
				AttributeType: v1.ConfigType_STRING,
				Description:   "the static password",
			},
		},
	}, nil
}

// Configure configures the plugin
func (p *Provider) Configure(ctx context.Context, req *v1.ConfigureRequest) (*empty.Empty, error) {
	if val, ok := req.Attributes[flagUsername]; ok {
		username, err := plugin.UnmarshalConfigurationValue(v1.ConfigType_STRING, val.GetValue())
		if err != nil {
			return nil, err
		}
		p.username = username.(string)
		log.Printf("configured username as '%s'", p.username)
	}
	if val, ok := req.Attributes[flagPassword]; ok {
		password, err := plugin.UnmarshalConfigurationValue(v1.ConfigType_STRING, val.GetValue())
		if err != nil {
			return nil, err
		}
		p.password = password.(string)
		log.Printf("configured password")
	}
	return &empty.Empty{}, nil
}
