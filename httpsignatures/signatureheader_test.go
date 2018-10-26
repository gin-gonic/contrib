package httpsignatures

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sampleKeyID            = "rsa-key-1"
	sampleAlgorithm        = "rsa-sha256"
	sampleSignature        = "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="
	sampleBodyContent      = "hello world"
	invalidSignatureString = "Signature keyId=\"rsa-key-1\",algorithm,headers=\"(request-target) host date digest \",signature=\"Hello world\""
)

var sampleHeader = []string{"(request-target)", "date", "digest"}
var sampleSignatureString = fmt.Sprintf("Signature keyId=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"",
	sampleKeyID, sampleAlgorithm, strings.Join(sampleHeader, " "), sampleSignature)

func newAuthorizationHeader(s string) http.Header {
	return http.Header{
		`Authorization`: []string{s},
	}
}

func newSignatureHeader(s string) http.Header {
	return http.Header{
		`Signature`: []string{s},
	}
}

func TestFromSignatureString(t *testing.T) {
	var tests = []struct {
		name      string
		header    http.Header
		keyID     string
		algorithm string
		headers   []string
		signature string
		err       error
	}{
		{
			name:   `empty headers`,
			header: http.Header{},
			err:    ErrNoSignature,
		},
		{
			name:   `Authorization Signature invalid begin`,
			header: newAuthorizationHeader(`notASignature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:    ErrInvalidAuthorizationHeader,
		},
		{
			name:   `Authorization Signature invalid key pair format`,
			header: newAuthorizationHeader(`Signature xxx`),
			err:    ErrUnterminatedParameter,
		},
		{
			name:   `Authorization Signature missing keyId`,
			header: newAuthorizationHeader(`Signature algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:    ErrMissingKeyID,
		},
		{
			name:   `Authorization Signature missing signature`,
			header: newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",headers=""`),
			err:    ErrMissingSignature,
		},
		{
			name:      `Authorization Signature header empty`,
			header:    newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"date"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Authorization Signature missing headers`,
			header:    newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"date"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Authorization Normal case`,
			header:    newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"(request-target)", "date", "digest"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Authorization Repeated params`,
			header:    newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",keyId="sample_key_id_2"`),
			err:       nil,
			keyID:     "sample_key_id_2",
			algorithm: "hmac-sha512",
			headers:   []string{"(request-target)", "date", "digest"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:   `Signature missing keyId`,
			header: newSignatureHeader(`algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:    ErrMissingKeyID,
		},
		{
			name:   `Signature missing signature`,
			header: newSignatureHeader(`keyId="sample_key_id",algorithm="hmac-sha512",headers=""`),
			err:    ErrMissingSignature,
		},
		{
			name:      `Signature header empty`,
			header:    newSignatureHeader(`keyId="sample_key_id",algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"date"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Signature missing headers`,
			header:    newSignatureHeader(`keyId="sample_key_id",algorithm="hmac-sha512",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"date"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Normal case`,
			header:    newSignatureHeader(`keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"(request-target)", "date", "digest"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Repeated params`,
			header:    newSignatureHeader(`keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",keyId="sample_key_id_2"`),
			err:       nil,
			keyID:     "sample_key_id_2",
			algorithm: "hmac-sha512",
			headers:   []string{"(request-target)", "date", "digest"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
	}

	for _, tc := range tests {

		r, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err, tc.name)
		r.Header = tc.header

		s, err := NewSignatureHeader(r)
		require.Equal(t, tc.err, err, tc.name)
		if err != nil {
			continue
		}
		assert.Equal(t, KeyID(tc.keyID), s.keyID, tc.name)
		assert.Equal(t, tc.algorithm, s.algorithm, tc.name)
		assert.Equal(t, tc.headers, s.headers, tc.name)
		assert.Equal(t, tc.signature, s.signature, tc.name)
	}
}
