package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/aws"
	v1 "github.com/VJftw/docker-registry-proxy/pkg/genproto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mozilla.org/pkcs7"
)

const (
	flagAvailabilityZones = "availability_zones"
	flagPrivateIPs        = "private_ips"
	flagInstanceIDs       = "instance_ids"
	flagAccountIDs        = "account_ids"
	flagImageIDs          = "image_ids"
	flagRegions           = "regions"
)

func main() {
	plugin.ServeAuthVerifierPlugin(NewVerifier())
}

// Verifier represents an AuthenticationVerifier
type Verifier struct {
	v1.AuthenticationVerifierServer
	v1.ConfigurationServer

	certs []*x509.Certificate

	// devpayProductCodes      []string
	// marketplaceProductCodes []string
	availabilityZones []string
	privateIPs        []string
	instanceIDs       []string
	// billingProducts         []string
	accountIDs []string
	imageIDs   []string
	// architectures           []string
	// kernelIDs               []string
	// ramdiskIDs              []string
	regions []string
}

// NewVerifier returns a new Verifier
func NewVerifier() *Verifier {
	certs, err := aws.GetCertificates(aws.AWSCertificates)
	if err != nil {
		log.Fatalf("could not get certificate pool: %s", err)
	}
	return &Verifier{
		certs: certs,
	}
}

// GetConfigurationSchema returns the schema for the plugin
func (v *Verifier) GetConfigurationSchema(ctx context.Context, _ *empty.Empty) (*v1.ConfigurationSchema, error) {
	return &v1.ConfigurationSchema{
		Attributes: map[string]*v1.ConfigurationAttribute{
			flagAvailabilityZones: &v1.ConfigurationAttribute{
				AttributeType: v1.ConfigType_STRING_SLICE,
				Description:   "the availability zones to accept",
			},
			flagPrivateIPs: &v1.ConfigurationAttribute{
				AttributeType: v1.ConfigType_STRING_SLICE,
				Description:   "the private IPs to accept",
			},
			flagInstanceIDs: &v1.ConfigurationAttribute{
				AttributeType: v1.ConfigType_STRING_SLICE,
				Description:   "the instance IDs to accept",
			},
			flagAccountIDs: &v1.ConfigurationAttribute{
				AttributeType: v1.ConfigType_STRING_SLICE,
				Description:   "the account IDs to accept",
			},
			flagImageIDs: &v1.ConfigurationAttribute{
				AttributeType: v1.ConfigType_STRING_SLICE,
				Description:   "the image IDs to accept",
			},
			flagRegions: &v1.ConfigurationAttribute{
				AttributeType: v1.ConfigType_STRING_SLICE,
				Description:   "the regions to accept",
			},
		},
	}, nil
}

// Configure configures the plugin
func (v *Verifier) Configure(ctx context.Context, req *v1.ConfigureRequest) (*empty.Empty, error) {
	v.availabilityZones = plugin.GetStringSliceValue(flagAvailabilityZones, req)
	v.privateIPs = plugin.GetStringSliceValue(flagPrivateIPs, req)
	v.instanceIDs = plugin.GetStringSliceValue(flagInstanceIDs, req)
	v.accountIDs = plugin.GetStringSliceValue(flagAccountIDs, req)
	v.imageIDs = plugin.GetStringSliceValue(flagImageIDs, req)
	v.regions = plugin.GetStringSliceValue(flagRegions, req)

	return &empty.Empty{}, nil
}

// Verify checks the given credentials
func (v *Verifier) Verify(ctx context.Context, req *v1.VerifyRequest) (*empty.Empty, error) {
	encodedPassword := req.GetPassword()
	instanceIdentityPassword := &aws.InstanceIdentityPassword{}
	if err := instanceIdentityPassword.Decode(encodedPassword); err != nil {
		return nil, fmt.Errorf("could not decode password: %w", err)
	}

	instanceIdentityPassword.Signature = []byte(fmt.Sprintf("-----BEGIN PKCS7-----\n%s\n-----END PKCS7-----", instanceIdentityPassword.Signature))

	decodedSig, _ := pem.Decode(instanceIdentityPassword.Signature)

	p7, err := pkcs7.Parse(decodedSig.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse PKCS7 signature '%s': %w", instanceIdentityPassword.Signature, err)
	}

	p7.Content = instanceIdentityPassword.Payload

	p7.Certificates = v.certs

	if err := p7.Verify(); err != nil {
		return nil, fmt.Errorf("could not verify signed data: %w", err)
	}

	doc := &aws.InstanceIdentityDocument{}
	if err := json.Unmarshal(p7.Content, doc); err != nil {
		return nil, fmt.Errorf("could not unmarshal document: %w", err)
	}

	if found := aws.CheckWhitelist(doc.AvailabilityZone, v.availabilityZones); !found {
		return nil, fmt.Errorf("%s: %w", doc.AvailabilityZone, aws.ErrNotWhitelisted)
	}

	if found := aws.CheckWhitelist(doc.PrivateIP, v.privateIPs); !found {
		return nil, fmt.Errorf("%s: %w", doc.PrivateIP, aws.ErrNotWhitelisted)
	}

	if found := aws.CheckWhitelist(doc.InstanceID, v.instanceIDs); !found {
		return nil, fmt.Errorf("%s: %w", doc.InstanceID, aws.ErrNotWhitelisted)
	}

	if found := aws.CheckWhitelist(doc.AccountID, v.accountIDs); !found {
		return nil, fmt.Errorf("%s: %w", doc.AccountID, aws.ErrNotWhitelisted)
	}

	if found := aws.CheckWhitelist(doc.ImageID, v.imageIDs); !found {
		return nil, fmt.Errorf("%s: %w", doc.ImageID, aws.ErrNotWhitelisted)
	}

	if found := aws.CheckWhitelist(doc.Region, v.regions); !found {
		return nil, fmt.Errorf("%s: %w", doc.Region, aws.ErrNotWhitelisted)
	}

	return &empty.Empty{}, nil
}
