// Package rx provide utility functions for channel inspired by RX: https://rxjs.dev/guide/operators
package rx

import (
	"sync"

	"go.uber.org/atomic"
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
	observer Observer[T],
	c <-chan T,
	errChs ...<-chan error,
) error {
	wg := new(sync.WaitGroup)
	wg.Add(1 + len(errChs))

	globalErr := new(atomic.Error)

	processErrs := func(errs <-chan error) {
		defer wg.Done()

		for err := range errs {
			if err != nil {
				if globalErr.Load() != nil {
					return
				}

				if observer.Err != nil {
					globalErr.Store(err)
					observer.Err(err)

					return
				}
			}
		}
	}

	for _, errs := range errChs {
		go processErrs(errs)
	}

	go func() {
		defer wg.Done()

		for v := range c {
			if globalErr.Load() != nil {
				break
			}

			if observer.Next != nil {
				observer.Next(v)
			}
		}
	}()

	wg.Wait()

	err := globalErr.Load()
	if err != nil {
		return err
	}

	if observer.Done != nil {
		observer.Done()
	}

	return nil
}

// Observable creates value channel and error channel and operate them follow observe function
func Observable[T any](observe func(observer Observer[T]), options ...Option) (<-chan T, <-chan error) {
	out, errs := observableChWithErrs[T](options...)
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

func observableChWithErrs[T any](options ...Option) (chan T, chan error) {
	return observableCh[T](options...), make(chan error)
}

func observableCh[T any](options ...Option) chan T {
	opts := parseOption(options...)

	return make(chan T, opts.bufferSize)
}
