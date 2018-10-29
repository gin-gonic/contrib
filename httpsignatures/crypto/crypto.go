package crypto

// Crypto interface for signing algorithim
type Crypto interface {
	Name() string
	Sign(msg string, secret string) ([]byte, error)
}
