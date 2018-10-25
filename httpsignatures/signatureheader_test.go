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

func TestFromSignatureString(t *testing.T) {
	var tests = []struct {
		// TODO: fill name of the test case
		name            string
		signatureHeader string
		success         bool
		keyID           string
		algorithm       string
		headers         []string
		signature       string
	}{
		{
			signatureHeader: ``,
			success:         false,
		},
		{
			signatureHeader: `Signature xxx`,
			success:         false,
		},
		// TODO fix this test case
		{
			signatureHeader: `notASignature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`,
			success:         false,
		},
		// TODO fix this test case
		{
			signatureHeader: `Signature algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`,
			success:         false,
		},
		// TODO fix this test case
		{
			signatureHeader: `Signature keyId="sample_key_id",algorithm="hmac-sha512",headers=""`,
			success:         false,
		},
		// TODO fix this test case
		{
			signatureHeader: `Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`,
			success:         true,
			keyID:           "sample_key_id",
			algorithm:       "hmac-sha512",
			headers:         []string{},
			signature:       "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			signatureHeader: `Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ=="`,
			success:         true,
			keyID:           "sample_key_id",
			algorithm:       "hmac-sha512",
			headers:         []string{"(request-target)", "date", "digest"},
			signature:       "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
		{
			signatureHeader: `Signature keyId="sample_key_id",algorithm="hmac-sha512",headers="(request-target) date digest",signature="70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",keyId="sample_key_id_2"`,
			success:         true,
			keyID:           "sample_key_id_2",
			algorithm:       "hmac-sha512",
			headers:         []string{"(request-target)", "date", "digest"},
			signature:       "70AaN3BDO0XC9QbtgksgCy2jJvmOvshq8VmjSthdXC+sgcgrKrl9WME4DbZv4W7UZKElvCemhDLHQ1Nln9GMkQ==",
		},
	}

	for _, tc := range tests {
		keyID, signature, headers, algorithm, err := fromSignatureString(tc.signatureHeader)
		require.True(t, (err == nil) == tc.success)
		assert.Equal(t, tc.keyID, keyID)
		assert.Equal(t, tc.algorithm, algorithm)
		assert.Equal(t, tc.headers, headers)
		assert.Equal(t, tc.signature, base64.StdEncoding.EncodeToString(signature))
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
