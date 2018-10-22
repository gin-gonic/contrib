package httpsignatures

import (
	"net/http"
)

// Validator interface for checking if a request is valid or not
type Validator interface {
	IsValid(*http.Request) bool
}
