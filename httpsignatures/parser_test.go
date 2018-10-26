package httpsignatures

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	var tests = []struct {
		input  string
		params map[string]string
	}{
		{
			input: `Signature keyId="rsa-key-1",algorithm="rsa-sha256",headers="(request-target) host date digest",signature="Hello world"`,
			params: map[string]string{
				"keyId":     "rsa-key-1",
				"algorithm": "rsa-sha256",
				"headers":   "(request-target) host date digest",
				"signature": "Hello world",
			},
		},
	}

	for _, tc := range tests {
		p, err := newParser(tc.input)
		require.NoError(t, err)
		results, err := p.parse()
		require.NoError(t, err)
		assert.Equal(t, tc.params, results)
	}
}
