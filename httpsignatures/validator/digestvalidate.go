package validator

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
)

//TODO: support more digest

//ErrInvalidDigest error when sha256 of body do not match with submitted digest
var ErrInvalidDigest = errors.New("Sha256 of body is not match with digest")

// DigestValidator checking digest in header match body
type DigestValidator struct {
}

// NewDigestValidator return pointer of new DigestValidator
func NewDigestValidator() *DigestValidator {
	return &DigestValidator{}
}

// Validate return error when checking digest match body
func (v *DigestValidator) Validate(r *http.Request) error {
	headerDigest := r.Header.Get("digest")
	digest, err := calculateDigest(r)
	if err != nil {
		return err
	}
	if digest != headerDigest {
		return ErrInvalidDigest
	}
	return nil
}

func calculateDigest(r *http.Request) (string, error) {
	if r.ContentLength == 0 {
		return "", nil
	}
	body, err := r.GetBody()
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, err = io.Copy(h, body)
	if err != nil {
		return "", err
	}
	digest := fmt.Sprintf("SHA-256=%s", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	return digest, nil
}
