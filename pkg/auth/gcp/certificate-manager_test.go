package gcp_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/gcp"
	"github.com/stretchr/testify/assert"
)

func TestGetMaxAgeFromHeader(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		in  string
		out time.Duration
	}{
		{"public, max-age=19933, must-revalidate, no-transform", 19933},
		{"public, max-age, must-revalidate, no-transform", 1800},
		{"public, max-age=123123, must-revalidate, no-transform", 123123},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			tt := tt
			t.Parallel()
			header := http.Header{}
			header.Set("cache-control", tt.in)
			age := gcp.GetMaxAgeFromHeader(header)
			assert.Equal(t, tt.out*time.Second, age)
		})
	}
}
