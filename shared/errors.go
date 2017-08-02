package shared

import "fmt"

// BadRequest returns an error that will be handled as an HTTP 400
func BadRequest(msg string) error {
	return fmt.Errorf("[BadRequest] %v", msg)
}

// InternalServerError returns an error that will be handled as an HTTP 500
func InternalServerError(msg string) error {
	return fmt.Errorf("[InternalServerError] %v", msg)
}
