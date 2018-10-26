package httpsignatures

import (
	"errors"
)

var (
	//ErrInvalidAuthorizationHeader error when get invalid format of Authorization header
	ErrInvalidAuthorizationHeader = errors.New("Authorization header format is incorrect")
	//ErrInvalidKeyID error when KeyID in header is not provided
	ErrInvalidKeyID = errors.New("Invalid keyId")
	//ErrDateNotInRange error when date not in aceptable range
	ErrDateNotInRange = errors.New("Date submit is not in aceptable range")
	//ErrIncorrectAlgorithm error when Algorithm in header is not match with secret key
	ErrIncorrectAlgorithm = errors.New("Algorithm is not match")
	//ErrHeaderNotEnough error when requiremts header do not appear on heder field
	ErrHeaderNotEnough = errors.New("Header feild is not match requirement")
	//ErrInvalidDigest error when sha256 of body do not match with submitted digest
	ErrInvalidDigest = errors.New("Sha256 of body is not match with digest")
	// ErrNoSignature error when no Signature not found in header
	ErrNoSignature = errors.New("No Signature header found in request")
	// ErrSignatureFormat error when signature format is invalid
	ErrSignatureFormat = errors.New("Invalid Signature header format")
	//ErrInvalidSign error when signing string do not match
	ErrInvalidSign = errors.New("Invalid sign")
	//ErrMissingKeyID error when keyId not in header
	ErrMissingKeyID = errors.New("keyId must be on header")
	//ErrMissingSignature error when signature not in header
	ErrMissingSignature = errors.New("signature must be on header")

	//ErrUnterminatedParameter err when could not parse value
	ErrUnterminatedParameter = errors.New("Unterminated parameter")
	//ErrMisingDoubleQuote err when after character = not have double quote
	ErrMisingDoubleQuote = errors.New(`Missing " after = character`)
	//ErrMisingEqualCharacter err when there is no character = before " or , character
	ErrMisingEqualCharacter = errors.New(`Missing = character =`)
)
