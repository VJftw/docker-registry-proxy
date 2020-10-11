package docker_test

import (
	"context"
	"testing"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/runtimes/docker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestGetRepository(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		in  string
		out string
	}{
		{"example.org/example:12345", "example.org/example"},
		{"example.org/example@sha256:12345", "example.org/example"},
		{"example.org/sub-dir/example@sha256:12345", "example.org/sub-dir/example"},
		{"example@sha256:12345", "index.docker.io/example"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			tt := tt
			t.Parallel()
			out, _ := docker.GetRepository(tt.in)
			assert.Equal(t, tt.out, out)
		})
	}
}

func TestResolveRepositoryAuthProvider(t *testing.T) {

	c1 := testProviderClient("test1")
	c2 := testProviderClient("test2")
	c3 := testProviderClient("test3")

	var authProviders = map[string]dockerregistryproxyv1.AuthenticationProviderAPIClient{
		"example.org":         c1,
		"example.org/bar":     c2,
		"example.org/foo/bar": c3,
	}
	var tests = []struct {
		in  string
		out dockerregistryproxyv1.AuthenticationProviderAPIClient
	}{
		{"example.org/example:12345", c1},
		{"example.org/bar:12345", c2},
		{"example.org/foo/bar:12345", c3},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			// tt := tt
			// t.Parallel()
			repo, err := docker.GetRepository(tt.in)
			assert.NoError(t, err)
			out, err := docker.ResolveRepositoryAuthProvider(repo, authProviders)
			assert.NoError(t, err)
			assert.Exactly(t, tt.out, out)
		})
	}
}

type testProviderClient string

func (c testProviderClient) Provide(ctx context.Context, in *dockerregistryproxyv1.ProvideRequest, opts ...grpc.CallOption) (*dockerregistryproxyv1.ProvideResponse, error) {
	return nil, nil
}
