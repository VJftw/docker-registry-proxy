package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/gcp"
	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
)

const (
	flagUsername = "username"
	flagAudience = "audience"
)

func main() {
	plugin.ServeAuthProviderPlugin(NewProvider())
}

// Provider represents an AuthenticationProvider
type Provider struct {
	dockerregistryproxyv1.AuthenticationProviderAPIServer
	dockerregistryproxyv1.ConfigurationAPIServer

	client   *http.Client
	username string
	audience string
}

// NewProvider returns a new Provider
func NewProvider() *Provider {
	return &Provider{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Provide returns credentials TODO: cache response from metadata in memory
func (p *Provider) Provide(ctx context.Context, req *dockerregistryproxyv1.ProvideRequest) (*dockerregistryproxyv1.ProvideResponse, error) {
	metaReq, err := http.NewRequest("GET", gcp.MetadataIdentity(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not construct metadata request: %w", err)
	}
	metaQuery := metaReq.URL.Query()
	metaQuery.Add("audience", p.audience)
	metaQuery.Add("format", "full")
	metaQuery.Add("licenses", "TRUE")
	metaReq.URL.RawQuery = metaQuery.Encode()
	metaReq.Header = *gcp.MetadataHeader

	metaResp, err := p.client.Do(metaReq)
	if err != nil {
		return nil, fmt.Errorf("could not execute metadata request: %w", err)
	}

	defer metaResp.Body.Close()
	jwtBytes, err := ioutil.ReadAll(metaResp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read metadata response: %w", err)
	}

	return &dockerregistryproxyv1.ProvideResponse{
		Username: p.username,
		Password: string(jwtBytes),
	}, nil
}

// GetConfigurationSchema returns the schema for the plugin
func (p *Provider) GetConfigurationSchema(ctx context.Context, _ *empty.Empty) (*dockerregistryproxyv1.GetConfigurationSchemaResponse, error) {
	return &dockerregistryproxyv1.GetConfigurationSchemaResponse{
		Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttribute{
			flagUsername: &dockerregistryproxyv1.ConfigurationAttribute{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
				Description:   "the routing username to provide credentials",
			},
			flagAudience: &dockerregistryproxyv1.ConfigurationAttribute{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
				Description:   "the aud in the signed JWT",
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
	if val, ok := req.Attributes[flagAudience]; ok {
		audience, err := plugin.UnmarshalConfigurationValue(dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING, val.GetValue())
		if err != nil {
			return nil, err
		}
		p.audience = audience.(string)
		log.Printf("configured audience as '%s'", p.audience)
	}
	return &empty.Empty{}, nil
}
