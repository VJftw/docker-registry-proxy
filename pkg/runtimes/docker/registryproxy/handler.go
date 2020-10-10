package registryproxy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// ProxyHandler returns a http.HandlerFunc
func ProxyHandler(opts *ProxyOpts) func(http.ResponseWriter, *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {

		fmt.Printf("received request: %#v\n", req)

		if req.URL.Path == "/_/v1/auth" {
			Authenticate(opts, res, req)
			return
		}

		reqAuth := GetAuthentication(req)
		if reqAuth == nil {
			UnauthenticatedResponse(res, req, opts.Upstream.Path)
			return
		}

		if req.URL.Path == "/v2/" {
			res.WriteHeader(http.StatusOK)
			return
		}

		repository := GetRepositoryFromPath(req.URL.Path)
		if repository != reqAuth.DockerRepository {
			// TODO: maybe use forbidden instead? check API docs
			UnauthenticatedResponse(res, req, opts.Upstream.Path)
			return
		}

		proxy, err := NewProxy(req, opts)
		if err != nil {
			logger.Warn("could not create proxy", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// add original path to request context
		newCtx := context.WithValue(req.Context(), CtxPath, req.URL.Path)
		req = req.WithContext(newCtx)
		proxy.ServeHTTP(res, req)
	}
}

// GetPathWithBase returns the path to send to the upstream registry
func GetPathWithBase(path, base string) string {
	if base != "" {
		parts := strings.Split(path, "/")
		parts = append(parts, "")
		copy(parts[3:], parts[2:])
		parts[2] = strings.TrimPrefix(base, "/")
		return strings.Join(parts, "/")
	}
	return path
}

// GetFullRepositoryFromPath returns the upstream Docker registry repository
func GetFullRepositoryFromPath(path string, base string) string {
	return GetRepositoryFromPath(GetPathWithBase(path, base))
}

// GetRepositoryFromPath returns the Docker registry repository
func GetRepositoryFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		return ""
	}
	parts = parts[2 : len(parts)-2]
	return strings.Join(parts, "/")
}

// GetImageFromPath returns the Docker image
func GetImageFromPath(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// GetScope returns the requested scope
func GetScope(req *http.Request, upstreamPath string) string {
	repository := GetFullRepositoryFromPath(req.URL.Path, upstreamPath)
	return fmt.Sprintf("repository:%s:pull", repository)
}

// GetRepositoryFromScope returns the repository from a scope
func GetRepositoryFromScope(scope string) string {
	parts := strings.Split(scope, ":")
	for i, p := range parts {
		if p == "repository" && len(parts) > i {
			return parts[i+1]
		}
	}
	return ""
}
