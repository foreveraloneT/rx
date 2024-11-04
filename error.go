package rx

import (
	"errors"
)

var (
	// ErrNotFound is an error type that indicates the result is not found in the channel
	ErrNotFound = errors.New("result not found")
)
