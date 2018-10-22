package httpsignatures

import (
	"crypto/hmac"
	"crypto/sha512"
)

const algoHmacSha512 = "hmac-sha512"

// HmacSha512 signing algorithm using hmac and sha512
type HmacSha512 struct {
}

func (h *HmacSha512) sign(msg string, secret string) ([]byte, error) {
	mac := hmac.New(sha512.New, []byte(secret))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func (h *HmacSha512) name() string {
	return algoHmacSha512
}
