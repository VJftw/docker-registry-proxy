package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/gcp"
	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
)

const (
	flagAudiences      = "audiences"
	flagProjectIDs     = "project_ids"
	flagProjectNumbers = "project_numbers"
	flagZones          = "zones"
	flagLicenseIDs     = "license_ids"
)

func main() {
	plugin.ServeAuthVerifierPlugin(NewVerifier())
}

// Verifier represents an AuthenticationVerifier
type Verifier struct {
	dockerregistryproxyv1.AuthenticationVerifierAPIServer
	dockerregistryproxyv1.ConfigurationAPIServer

	certificateManager *gcp.CertificateManager
	wg                 sync.WaitGroup

	audiences  []string
	projectIDs []string
	zones      []string
	licenseIDs []string
}

// NewVerifier returns a new Verifier
func NewVerifier() *Verifier {
	v := &Verifier{
		certificateManager: gcp.NewCertificateManager(),
	}
	go v.certificateManager.Run(&v.wg)
	return v
}

// GetConfigurationSchema returns the schema for the plugin
func (v *Verifier) GetConfigurationSchema(ctx context.Context, _ *empty.Empty) (*dockerregistryproxyv1.GetConfigurationSchemaResponse, error) {
	return &dockerregistryproxyv1.GetConfigurationSchemaResponse{
		Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttribute{
			flagAudiences: &dockerregistryproxyv1.ConfigurationAttribute{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING_SLICE,
				Description:   "the audiences to accept",
			},
			flagProjectIDs: &dockerregistryproxyv1.ConfigurationAttribute{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING_SLICE,
				Description:   "the project IDs to accept",
			},
			flagProjectNumbers: &dockerregistryproxyv1.ConfigurationAttribute{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING_SLICE,
				Description:   "the project numbers to accept",
			},
			flagZones: &dockerregistryproxyv1.ConfigurationAttribute{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING_SLICE,
				Description:   "the zones to accept",
			},
			flagLicenseIDs: &dockerregistryproxyv1.ConfigurationAttribute{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING_SLICE,
				Description:   "the license IDs to accept",
			},
		},
	}, nil
}

// Configure configures the plugin
func (v *Verifier) Configure(ctx context.Context, req *dockerregistryproxyv1.ConfigureRequest) (*empty.Empty, error) {
	v.audiences = plugin.GetStringSliceValue(flagAudiences, req)
	v.projectIDs = plugin.GetStringSliceValue(flagProjectIDs, req)
	v.zones = plugin.GetStringSliceValue(flagZones, req)
	v.licenseIDs = plugin.GetStringSliceValue(flagLicenseIDs, req)

	return &empty.Empty{}, nil
}

// Verify checks the given credentials
func (v *Verifier) Verify(ctx context.Context, req *dockerregistryproxyv1.VerifyRequest) (*empty.Empty, error) {
	claims, err := gcp.GetTokenClaims(req.GetPassword(), v.certificateManager)
	if err != nil {
		return nil, err
	}

	if found := gcp.CheckStringWhitelist(claims.Audience, v.audiences); !found {
		return nil, fmt.Errorf("%s: %w", claims.Audience, gcp.ErrAudienceNotWhitelisted)
	}
	if found := gcp.CheckStringWhitelist(claims.Google.ComputeEngine.ProjectID, v.projectIDs); !found {
		return nil, fmt.Errorf("%s: %w", claims.Google.ComputeEngine.ProjectID, gcp.ErrProjectIDNotWhitelisted)
	}
	if found := gcp.CheckStringWhitelist(claims.Google.ComputeEngine.Zone, v.zones); !found {
		return nil, fmt.Errorf("%s: %w", claims.Google.ComputeEngine.Zone, gcp.ErrZoneNotWhitelisted)
	}
	for _, licenseID := range claims.Google.ComputeEngine.LicenseIDs {
		if found := gcp.CheckStringWhitelist(licenseID, v.licenseIDs); !found {
			return nil, fmt.Errorf("%s: %w", licenseID, gcp.ErrLicenseIDWhitelisted)
		}
	}

	return &empty.Empty{}, nil
}
