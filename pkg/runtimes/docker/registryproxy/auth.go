package registryproxy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type bearerResponse struct {
	Token string `json:"token"`
	// ExpiresIn represents how long until the token expires in seconds (e.g. 3600)
	ExpiresIn uint64 `json:"expires_in"`
	// IssuedAt represents when the token was issued (e.g. 2009-11-10T23:00:00Z)
	IssuedAt time.Time `json:"issued_at"`
}

// Claims represents the claims for the generated token
type Claims struct {
	jwt.StandardClaims
	DockerRepository string `json:"dockerRepository"`
}

// Valid returns whether or not the claims are valid
func (c Claims) Valid() error {
	return c.StandardClaims.Valid()
}

var hmacSecret = []byte("CHANGEME") // TODO: make this a parameter

// GetAuthentication returns the authentication for a request
func GetAuthentication(req *http.Request) *Claims {
	authHeader := req.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(bearerToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSecret, nil
		})
		if err != nil {
			logger.Warn("invalid token", zap.Error(err))
			return nil
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims
		}
		logger.Warn("invalid token")
	}
	return nil
}

// UnauthenticatedResponse writes the unauthenticated response
func UnauthenticatedResponse(res http.ResponseWriter, req *http.Request, upstreamPath string) {
	// realm := "https://auth.docker.io/token"
	// service := "registry.docker.io"

	uri := url.URL{
		Scheme: "http",
		Host:   req.Host,
		Path:   "_/v1/auth",
	}
	if req.TLS != nil {
		uri.Scheme = "https"
	}
	if host := req.Header.Get("X-Forwarded-Host"); host != "" {
		uri.Host = host
	}
	if scheme := req.Header.Get("X-Forwarded-Proto"); scheme != "" {
		uri.Scheme = scheme
	}

	realm := uri.String()
	service := req.Host
	scope := GetScope(req, upstreamPath)
	res.Header().Set(
		"Www-Authenticate",
		fmt.Sprintf(
			`Bearer realm="%s",service="%s",scope="%s"`,
			realm,
			service,
			scope,
		),
	)
	res.WriteHeader(http.StatusUnauthorized)
}

// Authenticate handles authentication
func Authenticate(opts *ProxyOpts, res http.ResponseWriter, req *http.Request) {
	var repository string
	if len(opts.AuthVerifiers) > 0 {
		username, password, ok := req.BasicAuth()
		if !ok {
			logger.Warn("missing basic auth")
			UnauthenticatedResponse(res, req, opts.Upstream.Path)
			return
		}
		if _, ok := opts.AuthVerifiers[username]; !ok && len(opts.AuthVerifiers) > 0 {
			logger.Warn("unsupported user route", zap.String("route", username))
			UnauthenticatedResponse(res, req, opts.Upstream.Path)
			return
		}
		repository = GetRepositoryFromScope(req.URL.Query().Get("scope"))
		_, err := opts.AuthVerifiers[username].Verify(req.Context(), &dockerregistryproxyv1.VerifyRequest{
			Username:   username,
			Password:   password,
			Repository: repository,
		})
		if err != nil {
			logger.Warn("could not verify", zap.Error(err))
			UnauthenticatedResponse(res, req, opts.Upstream.Path)
			return
		}
	} else {
		repository = GetRepositoryFromScope(req.URL.Query().Get("scope"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		DockerRepository: repository,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
			NotBefore: time.Now().UTC().Unix(),
		},
	})

	// TODO: make the hmac an argument
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		logger.Warn("could not create token", zap.Error(err))
	}
	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(&bearerResponse{Token: tokenString}); err != nil {
		logger.Error("could not encode response", zap.Error(err))
	}
}

// MakeBasicAuth returns a basic auth string from a ProvideResponse
func MakeBasicAuth(username string, password string) string {

	basicAuthDecoded := fmt.Sprintf("%s:%s", username, password)

	return base64.StdEncoding.EncodeToString([]byte(basicAuthDecoded))
}
