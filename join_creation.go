package rx

import (
	"sync"
)

// Merge creates a channel that merges multiple channels into one by merging their emission
func Merge[T any](cs []<-chan T, options ...Option) <-chan T {
	out := observableCh[T](options...)

	wg := new(sync.WaitGroup)
	wg.Add(len(cs))

	output := func(c <-chan T) {
		defer wg.Done()

		for v := range c {
			out <- v
		}
	}

	for _, c := range cs {
		go output(c)
	}

	go func() {
		defer close(out)

		wg.Wait()
	}()

	return out
}
