package httpsignatures

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var defaultRequiredHeaders = []string{requestTarget, "date", "digest"}

// Authenticator is the gin authenticator middleware.
type Authenticator struct {
	secrets Secrects
	v       Validator
	headers []string
}

// Option is the option to the Authenticator constructor.
type Option func(*Authenticator)

// WithValidator configures the Authenticator to use custom validator.
// The default validator is timestamp based.
func WithValidator(v Validator) Option {
	return func(a *Authenticator) {
		a.v = v
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

	if a.v == nil {
		a.v = NewDateValidator()
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
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if !a.v.IsValid(c.Request) {
			c.AbortWithError(http.StatusBadRequest, ErrDateNotInRange)
			return
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
		signature, err := secret.Algorithm.sign(signString, secret.Key)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !bytes.Equal(signature, sigHeader.signature) {
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

	if secret.Algorithm.name() != algorithm {
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
		case "host":
			fieldValue = r.Host
		case "(request-target)":
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
