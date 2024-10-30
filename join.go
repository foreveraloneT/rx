package rx

import "sync"

// MergeAll merges all emitting channels into one by merging their emission
func MergeAll[T any](cs <-chan (<-chan T), options ...Option) <-chan T {
	out := observableCh[T](options...)

	wg := new(sync.WaitGroup)
	wg.Add(1) // for loop cs job

	output := func(c <-chan T) {
		defer wg.Done()

		for v := range c {
			out <- v
		}
	}

	go func() {
		defer wg.Done()

		for c := range cs {
			wg.Add(1)
			go output(c)
		}
	}()

	go func() {
		defer close(out)

		wg.Wait()
	}()

	return out
}

// SwitchAll merges all emitting channels into one by merging their emission.
// It will immediately switch to the current emitting rx.
func SwitchAll[T any](cs <-chan (<-chan T), options ...Option) <-chan T {
	out := observableCh[T](options...)

	wg := new(sync.WaitGroup)
	wg.Add(1) // for loop cs job

	process := func(c <-chan T) chan<- struct{} {
		terminate := make(chan struct{})

		go func() {
			defer wg.Done()

		LOOP:
			for {
				select {
				case v, ok := <-c:
					if !ok {
						break LOOP
					}

					out <- v
				case <-terminate:
					break LOOP
				}
			}
		}()

		return terminate
	}

	go func() {
		defer wg.Done()

		var terminate chan<- struct{}
		for c := range cs {
			if terminate != nil {
				terminate <- struct{}{}
			}

			wg.Add(1)
			terminate = process(c)
		}
	}()

	go func() {
		defer close(out)

		wg.Wait()
	}()

	return out
}
