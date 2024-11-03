package rx

import (
	"sync"
)

// MergeAll merges all emitting channels into one by merging their emission
func MergeAll[T any](sources <-chan Result[<-chan Result[T]], options ...Option) <-chan Result[T] {
	finalOptions := prepend(WithBufferSize(cap(sources)), options)

	results := resultCh[T](finalOptions...)
	tmp := resultCh[T](finalOptions...)

	wg := new(sync.WaitGroup)
	wg.Add(1) // for loop sources job

	process := func(source <-chan Result[T]) {
		defer wg.Done()

		for v := range source {
			tmp <- v

			if v.IsError() {
				return
			}
		}
	}

	go func() {
		defer wg.Done()

		for v := range sources {
			source, err := v.Get()
			if err != nil {
				tmp <- Err[T](err)

				return
			}

			wg.Add(1)
			go process(source)
		}
	}()

	go func() {
		defer close(tmp)

		wg.Wait()
	}()

	go func() {
		defer close(results)

		for v := range tmp {
			results <- v

			if v.IsError() {
				return
			}
		}
	}()

	return results
}

// SwitchAll merges all emitting channels into one by merging their emission.
// It will immediately switch to the current emitting channel.
func SwitchAll[T any](sources <-chan Result[<-chan Result[T]], options ...Option) <-chan Result[T] {
	finalOptions := prepend(WithBufferSize(cap(sources)), options)

	results := resultCh[T](finalOptions...)
	tmp := resultCh[T](finalOptions...)

	wg := new(sync.WaitGroup)
	wg.Add(1) // for loop sources job

	process := func(source <-chan Result[T]) chan<- struct{} {
		terminate := make(chan struct{})

		go func() {
			defer wg.Done()

			for {
				select {
				case <-terminate:
					return
				case v, ok := <-source:
					if !ok {
						return
					}

					tmp <- v

					if v.IsError() {
						return
					}
				}
			}
		}()

		return terminate
	}

	go func() {
		defer wg.Done()

		var terminate chan<- struct{}
		for v := range sources {
			if terminate != nil {
				close(terminate)
			}

			source, err := v.Get()
			if err != nil {
				tmp <- Err[T](err)

				return
			}

			wg.Add(1)
			terminate = process(source)
		}
	}()

	go func() {
		defer close(tmp)

		wg.Wait()
	}()

	go func() {
		defer close(results)

		for v := range tmp {
			results <- v

			if v.IsError() {
				return
			}
		}
	}()

	return results
}
