package rx

import (
	"time"
)

// Interval creates a channel that emits an integer at a fixed interval
func Interval(d time.Duration, options ...Option) <-chan int {
	out := observableCh[int](options...)

	go func() {
		defer close(out)

		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for i := 0; ; i++ {
			<-ticker.C
			out <- i
		}
	}()

	return out
}

// From creates a channel that emits the values from the provided slice
func From[T any](vv []T, options ...Option) <-chan T {
	out := observableCh[T](options...)

	go func() {
		defer close(out)

		for _, v := range vv {
			out <- v
		}
	}()

	return out
}

// Range creates a channel that emits a range of integers
func Range(start int, count int, options ...Option) <-chan int {
	out := observableCh[int](options...)

	go func() {
		defer close(out)

		for i := start; i < start+count; i++ {
			out <- i
		}
	}()

	return out
}
