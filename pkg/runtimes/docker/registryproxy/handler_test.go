package registryproxy_test

import (
	"testing"

	"github.com/VJftw/docker-registry-proxy/pkg/runtimes/docker/registryproxy"
	"gotest.tools/assert"
)

func TestGetPathWithBase(t *testing.T) {
	var tests = []struct {
		in   string
		base string
		out  string
	}{
		{"/v2/library/alpine/manifests/latest", "", "/v2/library/alpine/manifests/latest"},
		{"/v2/nginx/manifests/latest", "/library", "/v2/library/nginx/manifests/latest"},
		{"/v2/securego/gosec/manifests/latest", "", "/v2/securego/gosec/manifests/latest"},
		{"/v2/gosec/manifests/latest", "/securego", "/v2/securego/gosec/manifests/latest"},
		{"/v2/gosec/manifests/latest", "securego", "/v2/securego/gosec/manifests/latest"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			res := registryproxy.GetPathWithBase(tt.in, tt.base)
			assert.Equal(t, tt.out, res)
		})
	}
}

func TestGetFullRepositoryFromPath(t *testing.T) {
	var tests = []struct {
		in   string
		base string
		out  string
	}{
		{"/v2/library/alpine/manifests/latest", "", "library/alpine"},
		{"/v2/library/nginx/manifests/latest", "", "library/nginx"},
		{"/v2/nginx/manifests/latest", "library", "library/nginx"},
		{"/v2/securego/gosec/manifests/latest", "", "securego/gosec"},
		{"/v2/gosec/manifests/latest", "securego", "securego/gosec"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			res := registryproxy.GetFullRepositoryFromPath(tt.in, tt.base)
			assert.Equal(t, tt.out, res)
		})
	}
}

func TestGetRepositoryFromPath(t *testing.T) {
	var tests = []struct {
		in   string
		base string
		out  string
	}{
		{"/v2/library/alpine/manifests/latest", "", "library/alpine"},
		{"/v2/library/nginx/manifests/latest", "", "library/nginx"},
		{"/v2/nginx/manifests/latest", "library", "nginx"},
		{"/v2/securego/gosec/manifests/latest", "", "securego/gosec"},
		{"/v2/gosec/manifests/latest", "securego", "gosec"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			res := registryproxy.GetRepositoryFromPath(tt.in)
			assert.Equal(t, tt.out, res)
		})
	}
}

func TestGetImageFromPath(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{"/v2/library/alpine/manifests/latest", "latest"},
		{"/v2/library/nginx/manifests/latest", "latest"},
		{"/v2/nginx/manifests/latest", "latest"},
		{"/v2/securego/gosec/manifests/latest", "latest"},
		{"/v2/gosec/manifests/latest", "latest"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			res := registryproxy.GetImageFromPath(tt.in)
			assert.Equal(t, tt.out, res)
		})
	}
}
