package httpsignatures

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	readID             = KeyID("read")
	writeID            = KeyID("write")
	invaldAlgo         = "invalidAlgo"
	invalidSignature   = "Invalid Signature"
	requestNilBodySig  = "ewYjBILGshEmTDDMWLeBc9kQfIscSKxmFLnUBU/eXQCb0hrY1jh7U5SH41JmYowuA4p6+YPLcB9z/ay7OvG/Sg=="
	requestBodyContent = "hello world"
	requestBodyDigest  = "SHA-256=uU0nuZNNPgilLlLX2n2r+sSE7+N6U4DukIj3rOLvzek="
	requestBodySig     = "s8MEyer3dSpSsnL0+mQvUYgKm2S4AEX+hsvKmeNI7wgtLFplbCZtt8YOcySZrCyYbOJdPF1NASDHfupSuekecg=="
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
	requiredHeaders     = []string{"(request-target)", "date", "digest"}
	readOnlyPermissions = Permission{readID}
	allPermissions      = Permission{readID, writeID}
	dateValidator       = NewDateValidator()
	requestTime         = time.Date(2018, time.October, 22, 07, 00, 07, 00, time.UTC)
)

func runTest(secretKeys Secrects, headers []string, v Validator, permissions Permission, req *http.Request) *gin.Context {
	gin.SetMode(gin.TestMode)
	auth := NewAuthenticator(secretKeys, headers, v)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	auth.Authenticated(permissions)(c)
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
	c := runTest(secrets, requiredHeaders, NewDateValidator(), readOnlyPermissions, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrNoSignature)
}

func TestAuthenticatedHeaderInvalidSignature(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(authorizationHeader, "hello")
	c := runTest(secrets, requiredHeaders, dateValidator, readOnlyPermissions, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrSignatureFormat)
}

func TestAuthenticatedHeaderWrongPermission(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(writeID, algoHmacSha512, sampleHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Now().Format(http.TimeFormat))
	c := runTest(secrets, requiredHeaders, dateValidator, readOnlyPermissions, req)
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrNotEnoughPermission)
}

func TestAuthenticateDateNotAccept(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, sampleHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Date(1990, time.October, 20, 0, 0, 0, 0, time.UTC).Format(http.TimeFormat))
	c := runTest(secrets, requiredHeaders, dateValidator, readOnlyPermissions, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrDateNotInRange)
}

func TestAuthenticateInvalidRequiredHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	invalidRequiredHeaders := []string{"date"}
	sigHeader := generateSignature(readID, algoHmacSha512, invalidRequiredHeaders, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)

	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	c := runTest(secrets, requiredHeaders, dateValidator, readOnlyPermissions, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrHeaderNotEnough)
}

func TestAuthenticateInvalidAlgo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, invaldAlgo, sampleHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	c := runTest(secrets, requiredHeaders, dateValidator, readOnlyPermissions, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrIncorrectAlgorithm)
}

func TestInvalidSignNotBase64(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, sampleHeader, invalidSignature)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	c := runTest(secrets, requiredHeaders, dateValidator, readOnlyPermissions, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, base64.CorruptInputError(7))
}

func TestInvalidSign(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, sampleHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	c := runTest(secrets, requiredHeaders, dateValidator, readOnlyPermissions, req)
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
	auth := NewAuthenticator(secrets, requiredHeaders, &dateAlwaysValid{})
	r.Use(auth.Authenticated(readOnlyPermissions))
	r.GET("/", httpTestGet)

	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, sampleHeader, requestBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", requestTime.Format(http.TimeFormat))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestHttpValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	auth := NewAuthenticator(secrets, requiredHeaders, &dateAlwaysValid{})
	r.Use(auth.Authenticated(readOnlyPermissions))
	r.GET("/", httpTestGet)

	req, _ := http.NewRequest("GET", "/", nil)
	sigHeader := generateSignature(readID, algoHmacSha512, sampleHeader, requestNilBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", requestTime.Format(http.TimeFormat))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpValidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	auth := NewAuthenticator(secrets, requiredHeaders, &dateAlwaysValid{})
	r.Use(auth.Authenticated(readOnlyPermissions))
	r.POST("/", httpTestPost)

	req, _ := http.NewRequest("POST", "/", strings.NewReader(sampleBodyContent))
	sigHeader := generateSignature(readID, algoHmacSha512, sampleHeader, requestBodySig)
	req.Header.Set(authorizationHeader, sigHeader)
	req.Header.Set("Date", requestTime.Format(http.TimeFormat))
	req.Header.Set("Digest", requestBodyDigest)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body, _ := ioutil.ReadAll(w.Result().Body)
	assert.Equal(t, body, []byte(sampleBodyContent))
}
