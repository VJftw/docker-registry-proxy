package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/gcp"
	"github.com/golang/protobuf/ptypes/empty"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
)

const (
	flagUsername = "username"
)

func main() {
	plugin.ServeAuthProviderPlugin(NewProvider())
}

// Provider represents an AuthenticationProvider using GCP Instance Identity Documents
type Provider struct {
	client   *http.Client
	username string
}

// NewProvider returns a new Provider
func NewProvider() *Provider {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Provider{
		client: client,
	}
}

// Provide returns credentials TODO: cache response from metadata in memory
func (p *Provider) Provide(ctx context.Context, req *dockerregistryproxyv1.ProvideRequest) (*dockerregistryproxyv1.ProvideResponse, error) {
	metaReq, err := http.NewRequest("GET", gcp.MetadataToken(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not construct metadata request: %w", err)
	}
	metaReq.Header = *gcp.MetadataHeader

	metaResp, err := p.client.Do(metaReq)
	if err != nil {
		return nil, fmt.Errorf("could not execute metadata request: %w", err)
	}

	defer metaResp.Body.Close()
	tokenBytes, err := ioutil.ReadAll(metaResp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read metadata response: %w", err)
	}

	return &dockerregistryproxyv1.ProvideResponse{
		Username: p.username,
		Password: string(tokenBytes),
	}, nil
}

// GetConfigurationSchema returns the schema for the plugin
func (p *Provider) GetConfigurationSchema(ctx context.Context, _ *empty.Empty) (*dockerregistryproxyv1.ConfigurationSchema, error) {
	return &dockerregistryproxyv1.ConfigurationSchema{
		Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttribute{
			flagUsername: &dockerregistryproxyv1.ConfigurationAttribute{
				AttributeType: dockerregistryproxyv1.ConfigType_STRING,
				Description:   "the routing username to provide credentials",
			},
		},
	}, nil
}

// Configure configures the plugin
func (p *Provider) Configure(ctx context.Context, req *dockerregistryproxyv1.ConfigureRequest) (*empty.Empty, error) {
	if val, ok := req.Attributes[flagUsername]; ok {
		username, err := plugin.UnmarshalConfigurationValue(dockerregistryproxyv1.ConfigType_STRING, val.GetValue())
		if err != nil {
			return nil, err
		}
		p.username = username.(string)
		log.Printf("configured username as '%s'", p.username)
	}
	return &empty.Empty{}, nil
}
