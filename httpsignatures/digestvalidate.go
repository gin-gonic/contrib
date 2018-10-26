package httpsignatures

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

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
	if r.Body == nil {
		return "", nil
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if err == io.EOF {
			return "", nil
		}
		return "", err
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	h := sha256.New()
	h.Write(body)
	digest := fmt.Sprintf("SHA-256=%s", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	return digest, nil
}
