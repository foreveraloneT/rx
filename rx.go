// Package rx provide utility functions for channel inspired by RX: https://rxjs.dev/guide/operators
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
		if size >= 0 {
			c.bufferSize = size
		}
	}
}

func defaultConfig() *config {
	return &config{
		bufferSize: 0,
	}
}

func parseOption(opts ...Option) *config {
	c := defaultConfig()

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func resultCh[T any](opts ...Option) chan Result[T] {
	c := parseOption(opts...)

	return make(chan Result[T], c.bufferSize)
}

// Observer represents the observer with `Next`, `Err`, and `Done` handler functions
type Observer[T any] struct {
	Next func(T)
	Err  func(error)
	Done func()
}

// Observable creates an Result channel with the given observer function
func Observable[T any](observe func(observer Observer[T]), options ...Option) <-chan Result[T] {
	results := resultCh[T](options...)
	isDone := new(atomic.Bool)
	done := make(chan struct{})

	go func() {
		observer := Observer[T]{
			Next: func(value T) {
				if isDone.Load() {
					return
				}

				results <- Ok(value)
			},
			Err: func(err error) {
				if isDone.Load() {
					return
				}

				results <- Err[T](err)

				isDone.Store(true)
				close(done)
			},
			Done: func() {
				if isDone.Load() {
					return
				}

				isDone.Store(true)
				close(done)
			},
		}

		observe(observer)
	}()

	go func() {
		defer close(results)

		<-done
	}()

	return results
}

// Observe subscribes to the results channel and calls the observer handlers functions
func Observe[T any](results <-chan Result[T], observer Observer[T]) error {
	for result := range results {
		v, err := result.Get()
		if err != nil {
			if observer.Err != nil {
				observer.Err(err)
			}

			return err
		}

		if observer.Next != nil {
			observer.Next(v)
		}
	}

	if observer.Done != nil {
		observer.Done()
	}

	return nil
}

func prepend[T any](v T, slice []T) []T {
	out := make([]T, 0, len(slice)+1)
	out = append(out, v)
	out = append(out, slice...)

	return out
}
