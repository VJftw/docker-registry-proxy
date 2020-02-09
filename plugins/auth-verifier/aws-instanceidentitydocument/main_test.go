package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/aws"
	v1 "github.com/VJftw/docker-registry-proxy/pkg/genproto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"go.mozilla.org/pkcs7"
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

func TestVerify(t *testing.T) {
	documentResp := []byte(`{
    "devpayProductCodes" : null,
    "marketplaceProductCodes" : [ "1abc2defghijklm3nopqrs4tu" ],
    "availabilityZone" : "us-west-2b",
    "privateIp" : "10.158.112.84",
    "version" : "2017-09-30",
    "instanceId" : "i-1234567890abcdef0",
    "billingProducts" : null,
    "instanceType" : "t2.micro",
    "accountId" : "123456789012",
    "imageId" : "ami-5fb8c835",
    "pendingTime" : "2016-11-19T16:32:11Z",
    "architecture" : "x86_64",
    "kernelId" : null,
    "ramdiskId" : null,
    "region" : "us-west-2"
}`)
	toBeSigned, err := pkcs7.NewSignedData(documentResp)
	assert.NoError(t, err)
	key := genTestKey()
	pemCert := getCertificate(key)
	aws.AWSCertificates = map[string]string{
		"TEST": string(pemCert),
	}
	block, _ := pem.Decode([]byte(pemCert))
	cert, err := x509.ParseCertificate(block.Bytes)
	assert.NoError(t, err)
	err = toBeSigned.AddSigner(cert, key, pkcs7.SignerInfoConfig{})
	assert.NoError(t, err)
	toBeSigned.Detach()
	pkcs7Resp, err := toBeSigned.Finish()
	assert.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/latest/dynamic/instance-identity/document":
			w.Write(documentResp)
			w.WriteHeader(http.StatusOK)
		case "/latest/dynamic/instance-identity/pkcs7":
			w.Write(pkcs7Resp)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusOK)
		}
		if r.URL.Path == "/latest/api/token" {
			w.WriteHeader(http.StatusOK)
			return
		}
	}))
	aws.MetadataURL = ts.URL

	verifier := NewVerifier()
	marshalledUsername, err := plugin.MarshalConfigurationValue(
		v1.ConfigType_STRING,
		"_",
	)
	verifier.Configure(context.Background(), &v1.ConfigureRequest{
		Attributes: map[string]*v1.ConfigurationAttributeValue{
			"username": &v1.ConfigurationAttributeValue{
				AttributeType: v1.ConfigType_STRING,
				Value:         marshalledUsername,
			},
		},
	})

	instanceIdentityPassword := &aws.InstanceIdentityPassword{
		Payload:   documentResp,
		Signature: pkcs7Resp,
	}

	encodedInstanceIdentityPassword, err := instanceIdentityPassword.Encode()
	assert.NoError(t, err)

	_, err = verifier.Verify(context.Background(), &v1.VerifyRequest{
		Username: "_",
		Password: encodedInstanceIdentityPassword,
	})
	assert.NoError(t, err)

}
