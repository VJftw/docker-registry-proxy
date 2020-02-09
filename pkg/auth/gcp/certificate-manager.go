package gcp

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

// CertificateManager represents a certificate manager
type CertificateManager struct {
	// Certificates (kid:cert(*rsa.PublicKey))
	Certificates      map[string]*rsa.PublicKey
	certificatesMutex sync.Mutex

	closeCh chan struct{}
}

// RespCertificates represents the data structure returned by the endpoint (kid:cert)
type RespCertificates map[string]string

var rgxMaxAge = regexp.MustCompile("max-age=([0-9]*)")

// GoogleOAuth2CertsURL is the default endpoint for retrieving public keys to verify the given JWT with
var GoogleOAuth2CertsURL = "https://www.googleapis.com/oauth2/v1/certs"

// NewCertificateManager returns a new certificate manager
func NewCertificateManager() *CertificateManager {
	return &CertificateManager{
		closeCh:      make(chan struct{}),
		Certificates: map[string]*rsa.PublicKey{},
	}
}

// GetPublicKeyByKeyID returns the associated public key to given kid
func (cM *CertificateManager) GetPublicKeyByKeyID(kid string) (*rsa.PublicKey, error) {
	for i := 0; i < 10; i++ {
		cM.certificatesMutex.Lock()
		if len(cM.Certificates) > 0 {
			defer cM.certificatesMutex.Unlock()

			if cert, ok := cM.Certificates[kid]; ok {
				return cert, nil
			}
			return nil, fmt.Errorf("could not find certificate for kid: %s", kid)
		}
		cM.certificatesMutex.Unlock()
		time.Sleep(500 * time.Millisecond) // back-off
	}
	return nil, fmt.Errorf("timed out waiting for certificates to load")
}

// Run runs the certificate manager, this should most likely be executed as a go-routine
func (cM *CertificateManager) Run(wg *sync.WaitGroup) {
	logger.Info("running certificate manager")
	defer logger.Info("stopped certificate manager")
	defer wg.Done()

	for {
		cM.certificatesMutex.Lock()
		resp, err := http.Get(GoogleOAuth2CertsURL)
		if err != nil {
			logger.Error("could not retrieve Google OAuth2 certs", zap.Error(err))
		}

		var resCerts RespCertificates
		err = json.NewDecoder(resp.Body).Decode(&resCerts)
		if err != nil {
			logger.Error("could not decode Google OAuth2 certs", zap.Error(err))
		}

		for kid, encodedCert := range resCerts {
			block, _ := pem.Decode([]byte(encodedCert))
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				logger.Error("could not parse Google OAuth2 certs", zap.Error(err))
			} else {
				cM.Certificates[kid] = cert.PublicKey.(*rsa.PublicKey)
				logger.Info("stored certificate", zap.String("kid", kid))
			}
		}
		cM.certificatesMutex.Unlock()

		expiry := GetMaxAgeFromHeader(resp.Header)
		logger.Info("waiting to refresh certificates", zap.Duration("seconds", expiry))

		select {
		case <-cM.closeCh:
			return
		case <-time.After(expiry):
		}
	}
}

// Stop stops the certificate manager
func (cM *CertificateManager) Stop() {
	close(cM.closeCh)
}

// GetMaxAgeFromHeader returns the max-age value from the cache-control header, defaulting to 1800.
func GetMaxAgeFromHeader(header http.Header) time.Duration {
	maxAge := time.Duration(1800)
	cacheControl := header.Get("cache-control")
	if cacheControl != "" {
		matches := rgxMaxAge.FindStringSubmatch(cacheControl)
		if len(matches) > 1 {
			foundMaxAge, err := strconv.ParseInt(matches[1], 10, 64)
			if err == nil {
				maxAge = time.Duration(foundMaxAge)
			} else {
				logger.Error("could not parse int", zap.String("in", matches[1]), zap.Error(err))
			}
		}
	}

	return maxAge * time.Second
}
