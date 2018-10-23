package httpsignatures

// KeyID define type
type KeyID string

// Secret define secret key and algorithm that key use
type Secret struct {
	Key       string
	Algorithm Crypto
}

// Secrects map with keyID and secret
type Secrects map[KeyID]*Secret
