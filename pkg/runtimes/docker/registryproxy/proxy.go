package registryproxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	v1 "github.com/VJftw/docker-registry-proxy/pkg/genproto/v1"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
)

// These are the errors returned by this file of the pkg
var (
	ErrUpstreamCouldNotParseHost = errors.New("could not parse upstream host")
)

// ProxyOpts represents the options available for the ProxyHandler
type ProxyOpts struct {
	Upstream *url.URL

	Registry  name.Registry
	Auth      authn.Authenticator
	Transport http.RoundTripper

	AuthVerifiers map[string]v1.AuthenticationVerifierClient

	TokenHMAC []byte
}

// GetProxyOpts validates and returns the given proxy handler options
func GetProxyOpts(
	upstreamAddr string,
	upstreamAuthClient v1.AuthenticationProviderClient,
	authVerifiers map[string]v1.AuthenticationVerifierClient,
) (*ProxyOpts, error) {

	opts := &ProxyOpts{
		Transport: NewRedirectFollower(transport.NewRetry(http.DefaultTransport)),
	}

	upstreamURI, err := url.Parse(upstreamAddr)
	if err != nil {
		return nil, fmt.Errorf("could not parse upstream: %w", err)
	}
	if len(upstreamURI.Host) < 1 {
		return nil, ErrUpstreamCouldNotParseHost
	}
	opts.Upstream = upstreamURI

	registryOpts := []name.Option{}
	if upstreamURI.Scheme == "http" {
		registryOpts = append(registryOpts, name.Insecure)
	}

	r, err := name.NewRegistry(opts.Upstream.Host, registryOpts...)
	if err != nil {
		return nil, fmt.Errorf("could not parse registry: %w", err)
	}
	opts.Registry = r

	if upstreamAuthClient != nil {
		auth, err := RegistryAuthResolveGRPC(upstreamAuthClient)
		if err != nil {
			return nil, fmt.Errorf("could not resolve auth: %w", err)
		}
		opts.Auth = auth
	} else {
		defAuth, err := authn.DefaultKeychain.Resolve(r)
		if err != nil {
			return nil, fmt.Errorf("could not resolve default keychain auth: %w", err)
		}
		opts.Auth = defAuth
	}

	opts.AuthVerifiers = authVerifiers
	return opts, nil
}

// NewProxy returns a httputil.ReverseProxy for a request. This proxy should be used for subsequent requests of the same session.
func NewProxy(req *http.Request, opts *ProxyOpts) (*httputil.ReverseProxy, error) {
	proxy := httputil.NewSingleHostReverseProxy(opts.Upstream)

	// Create docker registry transport
	t, err := transport.New(
		opts.Registry,
		opts.Auth,
		opts.Transport,
		[]string{GetScope(req, opts.Upstream.Path)},
	)
	if err != nil {
		return nil, err
	}
	proxy.Transport = t

	proxy.Director = func(req *http.Request) {
		// Support TLS upstream from HTTP connections
		req.URL.Host = opts.Upstream.Host
		req.URL.Scheme = opts.Upstream.Scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = opts.Upstream.Host
		if opts.Upstream.Path != "" {
			pwb := GetPathWithBase(req.URL.Path, opts.Upstream.Path)
			req.URL.Path = pwb
		}
	}
	proxy.ModifyResponse = func(res *http.Response) error {
		originalPath, _ := res.Request.Context().Value(CtxPath).(string)
		fmt.Printf("%s: %d\n", originalPath, res.StatusCode)
		return nil
	}

	return proxy, nil
}
