package aws

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
)

// InstanceIdentityPassword represents the password struct to encode into the ProvideResponse
type InstanceIdentityPassword struct {
	Payload   []byte
	Signature []byte
}

// Encode returns a gob encoded representation of the Instance Identity Password
func (iip *InstanceIdentityPassword) Encode() (string, error) {
	var password bytes.Buffer
	enc := gob.NewEncoder(&password)
	err := enc.Encode(iip)
	if err != nil {
		return "", err
	}

	var compPassword bytes.Buffer
	comp := gzip.NewWriter(&compPassword)
	_, err = comp.Write(password.Bytes())
	if err != nil {
		return "", err
	}
	if err := comp.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(compPassword.Bytes()), nil
}

// Decode returns a gob encoded representation of the Instance Identity Password
func (iip *InstanceIdentityPassword) Decode(s string) error {

	compressedPassword, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("could not decode password: %w", err)
	}

	decomp, err := gzip.NewReader(bytes.NewBuffer(compressedPassword))
	if err != nil {
		return fmt.Errorf("could not create gzip decompressor: %w", err)
	}
	decodedPassword := &bytes.Buffer{}
	_, err = io.Copy(decodedPassword, decomp)
	if err != nil {
		return fmt.Errorf("could not copy decompressed bytes: :%w", err)
	}
	if err := decomp.Close(); err != nil {
		return fmt.Errorf("could not close decompressor: %w", err)
	}

	dec := gob.NewDecoder(bytes.NewBuffer(decodedPassword.Bytes()))
	err = dec.Decode(iip)
	if err != nil {
		return fmt.Errorf("could not decode gob password: %w", err)
	}

	return nil
}
