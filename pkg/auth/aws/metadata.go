package aws

import (
	"fmt"
	"net/http"
)

var (
	// MetadataURL is the endpoint to make requests to fetch instance metadata from
	MetadataURL = "http://169.254.169.254"
	Version     = func() string { return fmt.Sprintf("%s/latest", MetadataURL) }
	APIToken    = func() string { return fmt.Sprintf("%s/api/token", Version()) }

	MetadataIdentity          = func() string { return fmt.Sprintf("%s/dynamic/instance-identity/document", Version()) }
	MetadataIdentitySignature = func() string { return fmt.Sprintf("%s/dynamic/instance-identity/pkcs7", Version()) }
)

func GetMetadataHeader(token string) http.Header {
	return http.Header{
		"X-aws-ec2-metadata-token": []string{token},
	}
}

var TokenHeader = &http.Header{
	"X-aws-ec2-metadata-token-ttl-seconds": []string{"21600"},
}
