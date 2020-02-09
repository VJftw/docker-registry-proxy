package aws_test

import (
	"testing"

	"github.com/VJftw/docker-registry-proxy/pkg/auth/aws"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	in := &aws.InstanceIdentityPassword{
		Payload:   []byte("foo"),
		Signature: []byte("bar"),
	}

	encoded, err := in.Encode()
	assert.NoError(t, err)

	out := &aws.InstanceIdentityPassword{}
	err = out.Decode(encoded)
	assert.NoError(t, err)

	assert.Equal(t, in, out)
}
