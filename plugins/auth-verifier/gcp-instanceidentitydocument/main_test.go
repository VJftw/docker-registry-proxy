package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/gcp"
	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func genTestKey() *rsa.PrivateKey {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	return key
}

func getCertificate(key *rsa.PrivateKey) []byte {
	tml := x509.Certificate{
		NotBefore:    time.Now().UTC(),
		NotAfter:     time.Now().UTC().Add(10),
		SerialNumber: big.NewInt(123123),
		Subject: pkix.Name{
			CommonName:   "Test",
			Organization: []string{"Test"},
		},
		BasicConstraintsValid: true,
	}

	cert, _ := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
}

func getTokenStr(key *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		&gcp.Claims{
			ExpiredTime: gcp.ClaimsTime{Time: time.Now().UTC().Add(10 * time.Second)},
			Google: gcp.ClaimsGoogle{
				ComputeEngine: gcp.ClaimsGoogleComputeEngine{
					Zone: "europe-west2-b",
				},
			},
		},
	)
	token.Header["kid"] = "test"
	tokenStr, _ := token.SignedString(key)
	return tokenStr
}

func TestVerify(t *testing.T) {
	key := genTestKey()
	certBytes := getCertificate(key)
	serverResp, _ := json.Marshal(map[string]string{"test": string(certBytes)})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(serverResp); err != nil {
			log.Fatal(err)
		}
	}))
	defer ts.Close()
	gcp.GoogleOAuth2CertsURL = ts.URL

	verifier := NewVerifier()
	marshalledZones, err := plugin.MarshalConfigurationValue(
		dockerregistryproxyv1.ConfigType_STRING_SLICE,
		[]string{"europe-west2-b"},
	)
	assert.NoError(t, err)
	_, err = verifier.Configure(context.Background(), &dockerregistryproxyv1.ConfigureRequest{
		Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttributeValue{
			"zones": &dockerregistryproxyv1.ConfigurationAttributeValue{
				AttributeType: dockerregistryproxyv1.ConfigType_STRING_SLICE,
				Value:         marshalledZones,
			},
		},
	})
	assert.NoError(t, err)

	tokenStr := getTokenStr(key)
	log.Println(tokenStr)
	_, err = verifier.Verify(context.Background(), &dockerregistryproxyv1.VerifyRequest{
		Username:   "_",
		Password:   tokenStr,
		Host:       "",
		Repository: "",
		Image:      "",
	})
	assert.NoError(t, err)

}
