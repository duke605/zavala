package destiny2

import "fmt"

// SimpleError is returned for errors that have no metadata
type SimpleError string

func (e SimpleError) Error() string {
	return fmt.Sprintf("destiny2: %s", string(e))
}

const (
	// ErrWebAuthRequired is returned when an endpoint
	// is hit that requires auth credentials
	ErrWebAuthRequired SimpleError = "WebAuthRequired"

	// ErrUnautorized is returned when invalid/expired credentials are used
	ErrUnautorized SimpleError = "Unauthorized"

	// ErrUnknown is returned when the body of a response is not JSON and the
	// status code is not handled
	ErrUnknown SimpleError = "Unknown"

	// ErrNotFound is returned when a requested resource could not be found.
	// Also, weirdly, error can also be returned when an endpoint is hit with
	// a method the endpoint is not expecting
	ErrNotFound SimpleError = "NotFound"
)
