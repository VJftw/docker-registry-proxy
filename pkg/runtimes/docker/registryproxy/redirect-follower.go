package registryproxy

import (
	"net/http"
)

type ctxKey int

// context keys
const (
	CtxPath ctxKey = 0
)

// RedirectFollower represents a Transport RoundTripper that follows redirects
type RedirectFollower struct {
	http.RoundTripper
}

// NewRedirectFollower wraps another http.RoundTripper following redirect response codes
func NewRedirectFollower(wrappedRoundTripper http.RoundTripper) http.RoundTripper {
	return &RedirectFollower{
		RoundTripper: wrappedRoundTripper,
	}
}

// RoundTrip implements http.RoundTripper
func (t *RedirectFollower) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// follow redirects
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		l, err := resp.Location()
		if err != nil {
			return resp, err
		}
		// cContext := context.WithValue(req.Context(), CtxPath, req.URL.Path)
		cReq := req.Clone(req.Context())
		resp.Request = cReq
		cReq.Header = http.Header{}
		cReq.URL = l
		cReq.Host = l.Host
		cReq.RequestURI = l.RequestURI()
		cReq.Header.Set("X-Forwarded-Host", l.Host)
		resp.Request = cReq
		return resp, err
	}
	return resp, err
}
