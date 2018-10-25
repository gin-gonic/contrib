package httpsignatures

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/stretchr/testify/assert"
)

const (
	readID             = KeyID("read")
	writeID            = KeyID("write")
	invalidKeyID       = KeyID("invalid key")
	invaldAlgo         = "invalidAlgo"
	invalidSignature   = "Invalid Signature"
	requestNilBodySig  = "ewYjBILGshEmTDDMWLeBc9kQfIscSKxmFLnUBU/eXQCb0hrY1jh7U5SH41JmYowuA4p6+YPLcB9z/ay7OvG/Sg=="
	requestBodyContent = "hello world"
	requestBodyDigest  = "SHA-256=uU0nuZNNPgilLlLX2n2r+sSE7+N6U4DukIj3rOLvzek="
	requestBodySig     = "s8MEyer3dSpSsnL0+mQvUYgKm2S4AEX+hsvKmeNI7wgtLFplbCZtt8YOcySZrCyYbOJdPF1NASDHfupSuekecg=="
	requestHost        = "kyber.network"
	requestHostSig     = "+qpk6uAlILo/1YV1ZDK2suU46fbaRi5guOyg4b6aS4nWqLi9u57V6mVwQNh0s6OpfrVZwAYaWHCmQFCgJiZ6yg=="
)

var (
	hmacsha512 = &HmacSha512{}
	secrets    = Secrects{
		readID: &Secret{
			Key:       "1234",
			Algorithm: hmacsha512,
		},
		writeID: &Secret{
			Key:       "5678",
			Algorithm: hmacsha512,
		},
	}
	requiredHeaders = []string{"(request-target)", "date", "digest"}
	submitHeader    = []string{"(request-target)", "date", "digest"}
	submitHeader2   = []string{"(request-target)", "date", "digest", "host"}
	dateValidator   = NewDateValidator()
	requestTime     = time.Date(2018, time.October, 22, 07, 00, 07, 00, time.UTC)
)

func runTest(secretKeys Secrects, headers []string, v Validator, req *http.Request) *gin.Context {
	gin.SetMode(gin.TestMode)
	auth := NewAuthenticator(secretKeys, WithRequiredHeaders(headers), WithValidator(v))
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	auth.Authenticated()(c)
	return c
}

func generateSignature(keyID KeyID, algorithm string, headers []string, signature string) string {
	return fmt.Sprintf(
		"Signature keyId=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"",
		keyID, algorithm, strings.Join(headers, " "), signature,
	)
}

func TestAuthenticatedHeaderNoSignature(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	c := runTest(secrets, requiredHeaders, NewDateValidator(), req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrNoSignature)
}

func TestAuthenticatedHeaderInvalidSignature(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(authorizationHeader, "hello")
	c := runTest(secrets, requiredHeaders, dateValidator, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrSignatureFormat)
}

func TestAuthenticatedHeaderWrongKey(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(invalidKeyID, algoHmacSha512, submitHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	c := runTest(secrets, requiredHeaders, dateValidator, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrInvalidKeyID)
}

func TestAuthenticateDateNotAccept(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, submitHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Date(1990, time.October, 20, 0, 0, 0, 0, time.UTC).Format(http.TimeFormat))
	c := runTest(secrets, requiredHeaders, dateValidator, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrDateNotInRange)
}

func TestAuthenticateInvalidRequiredHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	invalidRequiredHeaders := []string{"date"}
	sigHeader := generateSignature(readID, algoHmacSha512, invalidRequiredHeaders, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)

	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	c := runTest(secrets, requiredHeaders, dateValidator, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrHeaderNotEnough)
}

func TestAuthenticateInvalidAlgo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, invaldAlgo, submitHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	c := runTest(secrets, requiredHeaders, dateValidator, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrIncorrectAlgorithm)
}

func TestInvalidSignNotBase64(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, submitHeader, invalidSignature)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	c := runTest(secrets, requiredHeaders, dateValidator, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, base64.CorruptInputError(7))
}

func TestInvalidSign(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, submitHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	c := runTest(secrets, requiredHeaders, dateValidator, req)
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrInvalidSign)
}

// mock interface always return true
type dateAlwaysValid struct{}

func (v *dateAlwaysValid) IsValid(r *http.Request) bool { return true }

func httpTestGet(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"success": true,
		})
}

func httpTestPost(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.Render(http.StatusOK, render.Data{Data: body})
}
func TestHttpInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	auth := NewAuthenticator(secrets, WithValidator(&dateAlwaysValid{}))
	r.Use(auth.Authenticated())
	r.GET("/", httpTestGet)

	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, submitHeader, requestBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", requestTime.Format(http.TimeFormat))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestHttpValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	auth := NewAuthenticator(secrets, WithValidator(&dateAlwaysValid{}))
	r.Use(auth.Authenticated())
	r.GET("/", httpTestGet)

	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, submitHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", requestTime.Format(http.TimeFormat))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpValidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	auth := NewAuthenticator(secrets, WithValidator(&dateAlwaysValid{}))
	r.Use(auth.Authenticated())
	r.POST("/", httpTestPost)

	req, err := http.NewRequest("POST", "/", strings.NewReader(sampleBodyContent))
	assert.NoError(t, err)
	sigHeader := generateSignature(readID, algoHmacSha512, submitHeader, requestBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", requestTime.Format(http.TimeFormat))
	req.Header.Set("Digest", requestBodyDigest)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, body, []byte(sampleBodyContent))
}

func TestHttpValidRequestHost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	auth := NewAuthenticator(secrets, WithValidator(&dateAlwaysValid{}))
	r.Use(auth.Authenticated())
	r.POST("/", httpTestPost)

	requestURL := fmt.Sprintf("http://%s/", requestHost)
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(sampleBodyContent))
	assert.NoError(t, err)
	sigHeader := generateSignature(readID, algoHmacSha512, submitHeader2, requestHostSig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", requestTime.Format(http.TimeFormat))
	req.Header.Set("Digest", requestBodyDigest)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, body, []byte(sampleBodyContent))
}
