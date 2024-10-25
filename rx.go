// package rx provide utility functions for channel inspired by RX: https://rxjs.dev/guide/operators
package rx

import (
	"sync/atomic"
)

type config struct {
	bufferSize int
}

// Option represents an option for the channel utility
type Option func(*config)

// WithBufferSize sets the buffer size of the channel
func WithBufferSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.bufferSize = size
		}
	}
}

// newConfig returns a new config with default values
func newConfig() *config {
	return &config{
		bufferSize: 0,
	}
}

func parseOption(opts ...Option) *config {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Observer represents the observer functions
type Observer[T any] struct {
	Next func(T)
	Err  func(error)
	Done func()
}

// Observe observes the values from the source channel and error channel and handles them using the provided observer functions
func Observe[T any](
	c <-chan T,
	errs <-chan error,
	observer Observer[T],
) {
LOOP:
	for {
		select {
		case v, ok := <-c:
			if !ok {
				if observer.Done != nil {
					observer.Done()
				}

				break LOOP
			}

			if observer.Next != nil {
				observer.Next(v)
			}
		case err := <-errs:
			if err != nil {
				if observer.Err != nil {
					observer.Err(err)
				}

				break LOOP
			}
		}
	}
}

// Observable creates value channel and error channel and operate them follow observe function
func Observable[T any](observe func(observer Observer[T]), options ...Option) (<-chan T, <-chan error) {
	opts := parseOption(options...)

	out := make(chan T, opts.bufferSize)
	errs := make(chan error, opts.bufferSize)
	done := make(chan struct{})

	isDone := new(atomic.Bool)

	go func() {
		observe(Observer[T]{
			Next: func(v T) {
				if isDone.Load() {
					return
				}

				out <- v
			},
			Err: func(err error) {
				if isDone.Load() {
					return
				}

				errs <- err

				isDone.Store(true)
				done <- struct{}{}
			},
			Done: func() {
				if isDone.Load() {
					return
				}

				isDone.Store(true)
				done <- struct{}{}
			},
		})
	}()

	go func() {
		defer close(out)
		defer close(errs)
		defer close(done)

		<-done
	}()

	return out, errs
}
