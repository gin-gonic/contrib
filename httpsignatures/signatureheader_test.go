package httpsignatures

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sampleKeyID            = "rsa-key-1"
	sampleAlgorithm        = "rsa-sha256"
	sampleSignature        = "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="
	sampleBodyContent      = "hello world"
	sampleDigest           = "SHA-256=uU0nuZNNPgilLlLX2n2r+sSE7+N6U4DukIj3rOLvzek="
	sampleDigestNil        = ""
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
			name:   `Signature invalid begin`,
			header: newAuthorizationHeader(`notASignature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:    ErrInvalidAuthorizationHeader,
		},
		{
			name:   `Signature invalid key pair format`,
			header: newAuthorizationHeader(`Signature xxx`),
			err:    ErrSignatureFormat,
		},
		{
			name:   `Signature missing keyId`,
			header: newAuthorizationHeader(`Signature algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:    ErrMissingKeyID,
		},
		{
			name:   `Signature missing signature`,
			header: newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",headers=""`),
			err:    ErrMissingSignature,
		},
		{
			name:   `Signature multiple space`,
			header: newAuthorizationHeader(`Signature     keyId="sample_key_id",algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:    ErrInvalidAuthorizationHeader,
		},
		{
			name:      `Signature header empty`,
			header:    newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"date"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Signature missing headers`,
			header:    newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"date"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Normal case`,
			header:    newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`),
			err:       nil,
			keyID:     "sample_key_id",
			algorithm: "hmac-sha512",
			headers:   []string{"(request-target)", "date", "digest"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			name:      `Repeated params`,
			header:    newAuthorizationHeader(`Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",keyId="sample_key_id_2"`),
			err:       nil,
			keyID:     "sample_key_id_2",
			algorithm: "hmac-sha512",
			headers:   []string{"(request-target)", "date", "digest"},
			signature: "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		//TODO: Test newSignatureHeader case
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
		assert.Equal(t, tc.keyID, s.keyID, tc.name)
		assert.Equal(t, tc.algorithm, s.algorithm, tc.name)
		assert.Equal(t, tc.headers, s.headers, tc.name)
		assert.Equal(t, tc.signature, base64.StdEncoding.EncodeToString(s.signature), tc.name)

	}
}

func TestCalculateDigest(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost", strings.NewReader(sampleBodyContent))
	digest, err := calculateDigest(r)
	assert.Equal(t, sampleDigest, digest)
	assert.Nil(t, err)
}

func TestCalculateDigestNilBody(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	digest, err := calculateDigest(r)
	assert.Equal(t, sampleDigestNil, digest)
	assert.Nil(t, err)
}

func TestNewSignatureHeaderNoSignature(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	s, err := NewSignatureHeader(r)
	assert.Nil(t, s)
	assert.Equal(t, err, ErrNoSignature)
}

func TestNewSignatureHeaderInvalidSignature(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	r.Header.Set(authorizationHeader, invalidSignatureString)
	s, err := NewSignatureHeader(r)
	assert.Nil(t, s)
	assert.Equal(t, err, ErrSignatureFormat)
}

func TestNewSignatureHeaderDateMissing(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	r.Header.Set(authorizationHeader, sampleSignatureString)
	s, err := NewSignatureHeader(r)
	assert.Nil(t, s)
	assert.NotNil(t, err)
}

func TestNewSignatureHeaderWithBodyWrongDigest(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost", strings.NewReader(sampleBodyContent))
	r.Header.Set(authorizationHeader, sampleSignatureString)
	currentTime := time.Now()
	r.Header.Set("Date", currentTime.Format(http.TimeFormat))
	r.Header.Set("Digest", "this is wrong digest")
	s, err := NewSignatureHeader(r)
	assert.Nil(t, s)
	assert.Equal(t, ErrInvalidDigest, err)
}

func TestNewSignatureHeaderNilBody(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	r.Header.Set(authorizationHeader, sampleSignatureString)
	currentTime := time.Now()
	r.Header.Set("Date", currentTime.Format(http.TimeFormat))
	s, err := NewSignatureHeader(r)
	assert.Nil(t, err)
	assert.Equal(t, KeyID(sampleKeyID), s.keyID)
	assert.Equal(t, sampleAlgorithm, s.algorithm)
	assert.Equal(t, sampleHeader, s.headers)
	assert.Equal(t, currentTime.Format(http.TimeFormat), s.date.Format(http.TimeFormat))
	assert.Equal(t, sampleDigestNil, s.digest)
}

func TestNewSignatureHeaderWithBody(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://localhost", strings.NewReader(sampleBodyContent))
	r.Header.Set(authorizationHeader, sampleSignatureString)
	currentTime := time.Now()
	r.Header.Set("Date", currentTime.Format(http.TimeFormat))
	r.Header.Set("Digest", sampleDigest)
	s, err := NewSignatureHeader(r)
	assert.Nil(t, err)
	assert.Equal(t, KeyID(sampleKeyID), s.keyID)
	assert.Equal(t, sampleAlgorithm, s.algorithm)
	assert.Equal(t, sampleHeader, s.headers)
	assert.Equal(t, currentTime.Format(http.TimeFormat), s.date.Format(http.TimeFormat))
	assert.Equal(t, sampleDigest, s.digest)
}
