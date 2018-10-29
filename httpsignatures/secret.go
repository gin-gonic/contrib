package httpsignatures

import "github.com/gin-gonic/contrib/httpsignatures/crypto"

// KeyID define type
type KeyID string

// Secret define secret key and algorithm that key use
type Secret struct {
	Key       string
	Algorithm crypto.Crypto
}

// Secrects map with keyID and secret
type Secrects map[KeyID]*Secret
