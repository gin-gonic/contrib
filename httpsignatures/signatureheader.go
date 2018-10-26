package httpsignatures

import (
	"net/http"
	"strings"
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
	signature string
	algorithm string
}

//NewSignatureHeader new instace of SignatureHeader
func NewSignatureHeader(r *http.Request) (*SignatureHeader, error) {
	keyID, signature, headers, algorithm, err := parseHTTPRequest(r)
	if err != nil {
		return nil, err
	}
	return &SignatureHeader{
		keyID:     KeyID(keyID),
		headers:   headers,
		signature: signature,
		algorithm: algorithm,
	}, nil
}

func parseHTTPRequest(r *http.Request) (keyID string, signature string, headers []string, algorithm string, err error) {
	s, err := getSignatureString(r)
	if err != nil {
		return keyID, signature, headers, algorithm, err
	}
	return parseSignatureString(s)
}

func parseSignatureString(s string) (keyID string, signature string, headers []string, algorithm string, err error) {
	p := newParser(s)
	results, err := p.parse()
	if err != nil {
		return keyID, signature, headers, algorithm, err
	}
	keyID, ok := results["keyId"]
	if !ok {
		return keyID, signature, headers, algorithm, ErrMissingKeyID
	}
	signature, ok = results["signature"]
	if !ok {
		return keyID, signature, headers, algorithm, ErrMissingSignature
	}
	headerString, ok := results["headers"]
	if !ok || len(headerString) == 0 {
		headers = []string{"date"}
	} else {
		headers = strings.Split(headerString, " ")
	}
	algorithm, _ = results["algorithm"]

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
