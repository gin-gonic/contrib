package httpsignatures

import (
	"errors"
)

var (
	//ErrInvalidAuthorizationHeader error when get invalid format of Authorization header
	ErrInvalidAuthorizationHeader = errors.New("Authorization header format is incorrect")
	//ErrInvalidKeyID error when KeyID in header is not provided
	ErrInvalidKeyID = errors.New("Invalid keyId")
	//ErrDateNotFound error when no date in header
	ErrDateNotFound = errors.New("There is no Date on Headers")
	//ErrIncorrectAlgorithm error when Algorithm in header is not match with secret key
	ErrIncorrectAlgorithm = errors.New("Algorithm is not match")
	//ErrHeaderNotEnough error when requiremts header do not appear on heder field
	ErrHeaderNotEnough = errors.New("Header feild is not match requirement")
	// ErrNoSignature error when no Signature not found in header
	ErrNoSignature = errors.New("No Signature header found in request")
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
