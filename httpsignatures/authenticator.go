package httpsignatures

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/httpsignatures/validator"
	"github.com/gin-gonic/gin"
)

const (
	requestTarget = "(request-target)"
	date          = "date"
	digest        = "digest"
	host          = "host"
)

var defaultRequiredHeaders = []string{requestTarget, date, digest}

// Authenticator is the gin authenticator middleware.
type Authenticator struct {
	secrets    Secrects
	validators []validator.Validator
	headers    []string
}

// Option is the option to the Authenticator constructor.
type Option func(*Authenticator)

// WithValidator configures the Authenticator to use custom validator.
// The default validator is timestamp based.
func WithValidator(validators []validator.Validator) Option {
	return func(a *Authenticator) {
		a.validators = validators
	}
}

// WithRequiredHeaders is list of all requires HTTP headers that the client
// have to include in the singing string for the request to be considered valid.
// If not provided, the created Authenticator instance will use defaultRequiredHeaders variable.
func WithRequiredHeaders(headers []string) Option {
	return func(a *Authenticator) {
		a.headers = headers
	}
}

// NewAuthenticator creates a new Authenticator instance with
// given allowed permissions and required header and secret keys.
func NewAuthenticator(secretKeys Secrects, options ...Option) *Authenticator {
	var a = &Authenticator{secrets: secretKeys}

	for _, fn := range options {
		fn(a)
	}

	if a.validators == nil {
		a.validators = []validator.Validator{
			validator.NewDateValidator(),
			validator.NewDigestValidator(),
		}
	}

	if len(a.headers) == 0 {
		a.headers = defaultRequiredHeaders
	}

	return a
}

// Authenticated returns a gin middleware which permits given permissions in parameter.
func (a *Authenticator) Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		sigHeader, err := NewSignatureHeader(c.Request)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		for _, v := range a.validators {
			if err := v.Validate(c.Request); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		}
		if !a.isValidHeader(sigHeader.headers) {
			c.AbortWithError(http.StatusBadRequest, ErrHeaderNotEnough)
			return
		}

		secret, err := a.getSecret(sigHeader.keyID, sigHeader.algorithm)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		signString := constructSignMessage(c.Request, sigHeader.headers)
		signature, err := secret.Algorithm.Sign(signString, secret.Key)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		signatureBase64 := base64.StdEncoding.EncodeToString(signature)
		if signatureBase64 != sigHeader.signature {
			c.AbortWithError(http.StatusUnauthorized, ErrInvalidSign)
			return
		}
		c.Next()
	}
}

// isValidHeader check if all server required header is in header list
func (a *Authenticator) isValidHeader(headers []string) bool {
	m := len(headers)
	for _, h := range a.headers {
		i := 0
		for i = 0; i < m; i++ {
			if h == headers[i] {
				break
			}
		}
		if i == m {
			return false
		}
	}
	return true
}

func (a *Authenticator) getSecret(keyID KeyID, algorithm string) (*Secret, error) {
	secret, ok := a.secrets[keyID]
	if !ok {
		return nil, ErrInvalidKeyID
	}

	if secret.Algorithm.Name() != algorithm {
		if algorithm != "" {
			return nil, ErrIncorrectAlgorithm
		}
	}
	return secret, nil
}

func constructSignMessage(r *http.Request, headers []string) string {
	var signBuffer bytes.Buffer
	for i, field := range headers {
		var fieldValue string
		switch field {
		case host:
			fieldValue = r.Host
		case requestTarget:
			fieldValue = fmt.Sprintf("%s %s", strings.ToLower(r.Method), r.URL.RequestURI())
		default:
			fieldValue = r.Header.Get(field)
		}
		signString := fmt.Sprintf("%s: %s", field, fieldValue)
		signBuffer.WriteString(signString)
		if i < len(headers)-1 {
			signBuffer.WriteString("\n")
		}
	}
	return signBuffer.String()
}
