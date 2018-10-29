package validator

import (
	"net/http"
)

// Validator interface for checking if a request is valid or not
type Validator interface {
	Validate(*http.Request) error
}
