package httpsignatures

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Authenticator is the gin authenticator middleware
type Authenticator struct {
	secrets Secrects
	v       Validator
	headers []string
}

// NewAuthenticator creates a new Authenticator instance with
// given allowed permissions and required header and secret keys
func NewAuthenticator(secretKeys Secrects, headers []string, v Validator) *Authenticator {
	return &Authenticator{secrets: secretKeys, headers: headers, v: v}
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
		var signString string
		switch field {
		case "digest":
			signString = fmt.Sprintf("%s: %s", field, r.Header.Get("Digest"))
		case "date":
			signString = fmt.Sprintf("%s: %s", field, r.Header.Get("Date"))
		case "(request-target)":
			signString = fmt.Sprintf("%s: %s %s", field, strings.ToLower(r.Method), r.URL.RequestURI())
		default:
			signString = fmt.Sprintf("%s: %s", field, r.Header.Get(field))
		}
		signBuffer.WriteString(signString)
		if i < len(headers)-1 {
			signBuffer.WriteString("\n")
		}
	}
	return signBuffer.String()
}
