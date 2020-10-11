package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/aws"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
)

const (
	flagUsername = "username"
	// HTTPTimeout is the amount of time to wait before a HTTP request should timeout.
	HTTPTimeout = 5 * time.Second
)

func main() {
	plugin.ServeAuthProviderPlugin(NewProvider())
}

// Provider represents an AuthenticationProvider using GCP Instance Identity Documents.
type Provider struct {
	client   *http.Client
	username string
}

// NewProvider returns a new Provider.
func NewProvider() *Provider {
	client := &http.Client{
		Timeout: HTTPTimeout,
	}
	return &Provider{
		client: client,
	}
}

// GetConfigurationSchema returns the schema for the plugin.
func (p *Provider) GetConfigurationSchema(_ context.Context, _ *empty.Empty) (
	*dockerregistryproxyv1.GetConfigurationSchemaResponse, error,
) {
	return &dockerregistryproxyv1.GetConfigurationSchemaResponse{
		Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttribute{
			flagUsername: {
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
				Description:   "the routing username to provide credentials",
			},
		},
	}, nil
}

// Configure configures the plugin.
func (p *Provider) Configure(
	_ context.Context,
	req *dockerregistryproxyv1.ConfigureRequest,
) (*empty.Empty, error) {
	if val, ok := req.Attributes[flagUsername]; ok {
		username, err := plugin.UnmarshalConfigurationValue(
			dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
			val.GetValue(),
		)
		if err != nil {
			return nil, err
		}
		p.username = username.(string)
		log.Printf("configured username as '%s'", p.username)
	}
	return &empty.Empty{}, nil
}

// Provide returns credentials TODO: cache response from metadata in memory.
func (p *Provider) Provide(ctx context.Context,
	_ *dockerregistryproxyv1.ProvideRequest,
) (*dockerregistryproxyv1.ProvideResponse, error) {

	// tokenReq, _ := http.NewRequest("PUT", aws.ApiToken(), nil)
	// tokenReq.Header = *aws.TokenHeader
	// tokenResp, err := p.client.Do(tokenReq)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not execute token request: %w", err)
	// }

	// token, err := ioutil.ReadAll(tokenResp.Body)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not read token from body: %w", err)
	// }

	metaReq, _ := http.NewRequestWithContext(ctx, "GET", aws.MetadataIdentity(), nil)
	// metaReq.Header = aws.GetMetadataHeader(string(token))

	metaResp, err := p.client.Do(metaReq)
	if err != nil {
		return nil, fmt.Errorf("could not execute metadata request: %w", err)
	}

	defer metaResp.Body.Close()
	metaJSONBytes, err := ioutil.ReadAll(metaResp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read metadata response: %w", err)
	}

	sigReq, _ := http.NewRequestWithContext(ctx, "GET", aws.MetadataIdentitySignature(), nil)
	// sigReq.Header = aws.GetMetadataHeader(string(token))

	sigResp, err := p.client.Do(sigReq)
	if err != nil {
		return nil, fmt.Errorf("could not execute metadata request: %w", err)
	}

	defer sigResp.Body.Close()
	sigBytes, err := ioutil.ReadAll(sigResp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read metadata response: %w", err)
	}

	instanceIdentityPassword := &aws.InstanceIdentityPassword{
		Payload:   metaJSONBytes,
		Signature: sigBytes,
	}

	encodedInstanceIdentityPassword, err := instanceIdentityPassword.Encode()
	if err != nil {
		return nil, fmt.Errorf("could not encode instance identity password: %w", err)
	}

	return &dockerregistryproxyv1.ProvideResponse{
		Username: p.username,
		Password: encodedInstanceIdentityPassword,
	}, nil
}
