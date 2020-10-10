package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/gcp"
	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestProvide(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("a")); err != nil {
			log.Fatal(err)
		}
	}))
	defer ts.Close()
	gcp.MetadataURL = ts.URL
	provider := NewProvider()
	marshalledUsername, err := plugin.MarshalConfigurationValue(
		dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
		"_test",
	)
	assert.NoError(t, err)
	marshalledAudience, err := plugin.MarshalConfigurationValue(
		dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
		"foo.example.org",
	)
	assert.NoError(t, err)

	_, err = provider.Configure(context.Background(), &dockerregistryproxyv1.ConfigureRequest{
		Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttributeValue{
			"username": &dockerregistryproxyv1.ConfigurationAttributeValue{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
				Value:         marshalledUsername,
			},
			"audience": &dockerregistryproxyv1.ConfigurationAttributeValue{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
				Value:         marshalledAudience,
			},
		},
	})
	assert.NoError(t, err)

	resp, err := provider.Provide(context.Background(), &dockerregistryproxyv1.ProvideRequest{})
	assert.NoError(t, err)

	assert.Equal(t, "_test", resp.GetUsername())
	assert.Equal(t, "a", resp.GetPassword())
}
