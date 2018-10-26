package httpsignatures

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	authorizationHeader           = "Authorization"
	authorizationHeaderInitString = "Signature "
	signatureHeader               = "Signature"
	requestTarget                 = "(request-target)"
)

//SignatureHeader contains basic info signature header
type SignatureHeader struct {
	keyID     KeyID
	headers   []string
	signature []byte
	algorithm string
	date      time.Time
	digest    string
}

//NewSignatureHeader new instace of SignatureHeader
func NewSignatureHeader(r *http.Request) (*SignatureHeader, error) {
	keyID, signature, headers, algorithm, err := fromHTTPRequest(r)
	if err != nil {
		return nil, err
	}

	date, err := http.ParseTime(r.Header.Get("Date"))
	if err != nil {
		return nil, err
	}
	digest := r.Header.Get("Digest")
	for _, h := range headers {
		if h == "digest" {
			httpDigest, err := calculateDigest(r)
			if err != nil {
				return nil, err
			}
			if digest != httpDigest {
				return nil, ErrInvalidDigest
			}
			break
		}
	}
	return &SignatureHeader{
		keyID:     KeyID(keyID),
		headers:   headers,
		date:      date,
		signature: signature,
		algorithm: algorithm,
		digest:    digest,
	}, nil
}

func fromHTTPRequest(r *http.Request) (keyID string, signature []byte, headers []string, algorithm string, err error) {
	s, err := getSignatureString(r)
	if err != nil {
		return keyID, signature, headers, algorithm, err
	}
	return fromSignatureString(s)
}

func fromSignatureString(s string) (keyID string, signature []byte, headers []string, algorithm string, err error) {
	sigStructs := strings.Split(s, ",")
	for _, pair := range sigStructs {
		key, val, err := keyValSplit(pair)
		if err != nil {
			return keyID, signature, headers, algorithm, err
		}
		switch key {
		case "keyId":
			keyID = val
		case "signature":
			data, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return keyID, signature, headers, algorithm, err
			}
			signature = data
		case "headers":
			headers = strings.Split(val, " ")
		case "algorithm":
			algorithm = val
		}
	}
	return keyID, signature, headers, algorithm, nil
}

func getSignatureString(r *http.Request) (string, error) {
	if s, ok := r.Header[authorizationHeader]; ok {
		if strings.Index(s[0], authorizationHeaderInitString) != 0 {
			return "", ErrInvalidAuthorizationHeader
		}
		return strings.TrimPrefix(s[0], authorizationHeaderInitString), nil
	} else if s, ok = r.Header[signatureHeader]; ok {
		return s[0], nil
	}
	return "", ErrNoSignature
}

func keyValSplit(s string) (key string, val string, err error) {
	stringList := strings.SplitN(s, "=", 2)
	if len(stringList) < 2 {
		return "", "", ErrSignatureFormat
	}
	key = stringList[0]
	val = strings.Trim(stringList[1], "\"")
	return key, val, nil
}

func calculateDigest(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	h := sha256.New()
	h.Write(body)
	digest := fmt.Sprintf("SHA-256=%s", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	return digest, nil
}
