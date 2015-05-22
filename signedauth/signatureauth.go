package signedauth

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"hash"
	"net/http"
	"strings"
)

// SignedAuthManager defines the functions needed to fulfill an auth key managing role.
type SignedAuthManager interface {
	AuthHeaderPrefix() string                         // The beginning of the string from the HTTP_AUTHORIZATION header. The exact header must be followed by a space.
	SecretKey(string, *http.Request) (string, *Error) // The secret key for the provided access key and request. Header verification should happen here, and an error returned to fail.
	DataToSign(*http.Request) (string, *Error)        // The data which must be signed and verified, or an error to return.
	AuthHeaderRequired() bool                         // Whether or not a request without any header should be accepted (c.Next) or forbidden (c.Fail with status 403).
	HashFunction() func() hash.Hash                   // Returns the hash function to use, e.g. sha1.New (imported from "crypto/sha1"), or sha512.New384 for SHA-384.
	ContextKey() string                               // The key in the context where will be set the appropriate value if the request was correctly signed.
	ContextValue(string) interface{}                  // The value which will be stored in the context if authentication is successful, from the access key.
}

// Error defines the authentication failure with a status. The error string will *not* be returned by Gin.
type Error struct {
	Status int   // The status for this failure.
	Err    error // The error associated to this failure.
}

// SignatureAuth is the middleware function. It must be called with a struct which implements the SignedAuthManager interface.
func SignatureAuth(mgr SignedAuthManager) gin.HandlerFunc {

	return func(c *gin.Context) {
		accesskey, signature, err := extractAuthInfo(mgr, c.Request.Header.Get("Authorization"))
		if err != nil {
			// Credentials doesn't match, we return 401 Unauthorized and abort request.
			c.Fail(err.Status, err.Err)
		} else if accesskey == "" && signature == "" && !mgr.AuthHeaderRequired() {
			c.Next()
		} else {
			// Authorization header has the correct format.
			secret, keyerr := mgr.SecretKey(accesskey, c.Request)
			if keyerr != nil {
				c.Fail(keyerr.Status, keyerr.Err)
			} else {
				data, dataerr := mgr.DataToSign(c.Request)
				if dataerr != nil {
					c.Fail(dataerr.Status, dataerr.Err)
				} else if !isSignatureValid(mgr.HashFunction(), secret, data, signature) {
					// Accesskey is valid but signature is not.
					c.Fail(http.StatusUnauthorized, errors.New("Wrong access key or signature."))
				} else {
					// Accesskey and signature are valid.
					c.Set(mgr.ContextKey(), mgr.ContextValue(accesskey))
					c.Next()
				}
			}
		}
	}
}

// extractAuthInfo extracts the authentication information from the provided auth string.
func extractAuthInfo(mgr SignedAuthManager, auth string) (string, string, *Error) {
	if strings.HasPrefix(auth, mgr.AuthHeaderPrefix()+" ") {
		splitheader := strings.Split(auth, " ")
		if len(splitheader) != 2 {
			return "", "", &Error{http.StatusUnauthorized, errors.New("Invalid authorization header.")}
		}

		splitauth := strings.Split(splitheader[1], ":")
		if len(splitauth) != 2 {
			return "", "", &Error{http.StatusUnauthorized, errors.New("Invalid format for access key and signature.")}
		}
		return splitauth[0], splitauth[1], nil

	} else if mgr.AuthHeaderRequired() {
		return "", "", &Error{http.StatusUnauthorized, errors.New("Invalid authorization header.")}
	} else {
		return "", "", nil
	}

}

// isSignatureValid signs the request with the provided secret, and returns that signature.
func isSignatureValid(hashFunc func() hash.Hash, secret string, data string, signature string) bool {
	hash := hmac.New(hashFunc, []byte(secret))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)) == signature
}
