package gcp

import (
	"fmt"
	"net/http"
)

var (
	// MetadataURL is the endpoint to make requests to fetch instance metadata from
	MetadataURL           = "http://metadata.google.internal./computeMetadata/v1"
	serviceAccounts       = func() string { return fmt.Sprintf("%s/instance/service-accounts", MetadataURL) }
	defaultServiceAccount = "default"

	MetadataIdentity = func() string { return fmt.Sprintf("%s/%s/identity", serviceAccounts(), defaultServiceAccount) }
	MetadataScopes   = func() string { return fmt.Sprintf("%s/%s/scopes", serviceAccounts(), defaultServiceAccount) }
	MetadataToken    = func() string { return fmt.Sprintf("%s/%s/token", serviceAccounts(), defaultServiceAccount) }
	MetadataEmail    = func() string { return fmt.Sprintf("%s/%s/email", serviceAccounts(), defaultServiceAccount) }
)

var MetadataHeader = &http.Header{
	"Metadata-Flavor": []string{"Google"},
}
