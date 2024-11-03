package rx

// Result represents value emitted by an Observable
type Result[T any] struct {
	value T
	err   error
}

// Get returns value or error
func (r Result[T]) Get() (T, error) {
	if r.err != nil {
		return empty[T](), r.err
	}

	return r.value, nil
}

// IsError returns true if there is an error, otherwise false
func (r Result[T]) IsError() bool {
	return r.err != nil
}

// Value returns value if no error, otherwise it returns zero value
func (r Result[T]) Value() T {
	if r.err != nil {
		return empty[T]()
	}

	return r.value
}

// Error returns error
func (r Result[T]) Error() error {
	return r.err
}

// Ok returns Result when a value is valid
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err returns Result when value is invalid
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func empty[T any]() T {
	var empty T

	return empty
}
