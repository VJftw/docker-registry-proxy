package gcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Errors returned by Verifier
var (
	ErrAudienceNotWhitelisted      = errors.New("audience not whitelisted")
	ErrProjectIDNotWhitelisted     = errors.New("project id not whitelisted")
	ErrProjectNumberNotWhitelisted = errors.New("project number not whitelisted")
	ErrZoneNotWhitelisted          = errors.New("zone not whitelisted")
	ErrLicenseIDWhitelisted        = errors.New("license id not whitelisted")

	ErrTokenClaimsParse   = errors.New("unable to parse claims in token")
	ErrTokenClaimsExpired = errors.New("token claims expired")
	ErrTokenMissingKeyID  = errors.New("missing kid in token")
)

func CheckStringWhitelist(value string, whitelist []string) bool {
	if len(whitelist) < 1 {
		return true
	}

	for _, i := range whitelist {
		if value == i {
			return true
		}
	}

	return false
}

func CheckUint64Whitelist(value uint64, whitelist []uint64) bool {
	if len(whitelist) < 1 {
		return true
	}

	for _, i := range whitelist {
		if value == i {
			return true
		}
	}

	return false
}

func GetTokenClaims(tokenStr string, certManager *CertificateManager) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Header["kid"]; !ok {
			return nil, ErrTokenMissingKeyID
		}
		kid := token.Header["kid"].(string)
		pubKey, err := certManager.GetPublicKeyByKeyID(kid)
		if err != nil {
			return nil, fmt.Errorf("invalid kid in token: %w", err)
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrTokenClaimsParse
	}

	if err := claims.Valid(); err != nil {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}

	return claims, nil
}

type ClaimsTime struct {
	time.Time
}

func (t *ClaimsTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UnixNano())
}

func (t *ClaimsTime) UnmarshalJSON(b []byte) error {
	unix, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(unix, 0)
	return nil
}

// Claims represents the expected claims of the JWT returned by GCE instance metadata.
// https://cloud.google.com/compute/docs/instances/verifying-instance-identity#token_format
type Claims struct {
	TokenIssuer     string       `json:"iss"`
	IssuedTime      ClaimsTime   `json:"iat"`
	ExpiredTime     ClaimsTime   `json:"exp"`
	Audience        string       `json:"aud"`
	Subject         string       `json:"sub"`
	AuthorizedParty string       `json:"azp"`
	Google          ClaimsGoogle `json:"google"`
	Email           string       `json:"email"`
	EmailVerified   bool         `json:"email_verified"`
}

// ClaimsGoogle represents the Google subclaims
type ClaimsGoogle struct {
	ComputeEngine ClaimsGoogleComputeEngine `json:"compute_engine"`
}

// ClaimsGoogleComputeEngine represents the Google Compute Engine subclaims
type ClaimsGoogleComputeEngine struct {
	ProjectID         string     `json:"project_id"`
	ProjectNumber     uint64     `json:"project_number"`
	Zone              string     `json:"zone"`
	InstanceID        string     `json:"instance_id"`
	InstanceName      string     `json:"instance_name"`
	CreationTimestamp ClaimsTime `json:"instance_creation_timestamp"`
	LicenseIDs        []string   `json:"license_id"`
}

// Valid validates that the given claims are valid
func (c *Claims) Valid() error {
	if c.ExpiredTime.Sub(time.Now().UTC()) < 0 {
		return ErrTokenClaimsExpired
	}
	return nil
}
