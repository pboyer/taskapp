package shared

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// BadRequest returns an error that will be handled as an HTTP 400
func BadRequest(input interface{}) error {
	// Errors returned to API Gateway must be free of whitespace
	return fmt.Errorf(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, fmt.Sprintf("[BadRequest] %v", input)))
}

// InternalServerError returns an error that will be handled as an HTTP 500
func InternalServerError(msg interface{}) error {
	// Apex requires logging to be done to stdErr.
	fmt.Fprintf(os.Stderr, "[InternalServerError] %v", msg)

	// Obfuscate internal server errors
	return errors.New("[InternalServerError]")
}
