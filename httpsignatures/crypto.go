package httpsignatures

// Crypto interface for signing algorithim
type Crypto interface {
	name() string
	sign(msg string, secret string) ([]byte, error)
}
