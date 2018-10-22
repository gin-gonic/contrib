package httpsignatures

import (
	"errors"
)

var (
	//ErrInvalidKeyID error when KeyID in header is not provided
	ErrInvalidKeyID = errors.New("Invalid KeyID")
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
	//ErrNotEnoughPermission error when keyid do not have enough permission
	ErrNotEnoughPermission = errors.New("KeyID do not have permission")
	//ErrInvalidSign error when signing string do not match
	ErrInvalidSign = errors.New("Invalid sign")
)
