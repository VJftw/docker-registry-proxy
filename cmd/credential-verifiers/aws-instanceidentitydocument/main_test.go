package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/aws"
	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
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
	backupAWSCerts := aws.AWSCertificates
	aws.AWSCertificates = map[string]string{
		"TEST": string(pemCert),
	}
	defer func() {
		aws.AWSCertificates = backupAWSCerts
	}()
	block, _ := pem.Decode([]byte(pemCert))
	cert, err := x509.ParseCertificate(block.Bytes)
	assert.NoError(t, err)
	err = toBeSigned.AddSigner(cert, key, pkcs7.SignerInfoConfig{})
	assert.NoError(t, err)
	toBeSigned.Detach()
	signed, err := toBeSigned.Finish()
	pkcs7Resp := pem.EncodeToMemory(&pem.Block{Type: "PKCS7", Bytes: signed})
	pkcs7Resp = bytes.ReplaceAll(pkcs7Resp, []byte("-----BEGIN PKCS7-----"), []byte(""))
	pkcs7Resp = bytes.ReplaceAll(pkcs7Resp, []byte("-----END PKCS7-----"), []byte(""))
	// assert.Equal(t, 1, string(pkcs7Resp))
	assert.NoError(t, err)

	verifier := NewVerifier()
	marshalledUsername, err := plugin.MarshalConfigurationValue(
		dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
		"_",
	)
	assert.NoError(t, err)
	_, err = verifier.Configure(context.Background(), &dockerregistryproxyv1.ConfigureRequest{
		Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttributeValue{
			"username": &dockerregistryproxyv1.ConfigurationAttributeValue{
				AttributeType: dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING,
				Value:         marshalledUsername,
			},
		},
	})
	assert.NoError(t, err)

	instanceIdentityPassword := &aws.InstanceIdentityPassword{
		Payload:   documentResp,
		Signature: pkcs7Resp,
	}

	encodedInstanceIdentityPassword, err := instanceIdentityPassword.Encode()
	assert.NoError(t, err)

	_, err = verifier.Verify(context.Background(), &dockerregistryproxyv1.VerifyRequest{
		Username: "_",
		Password: encodedInstanceIdentityPassword,
	})
	assert.NoError(t, err)

}
