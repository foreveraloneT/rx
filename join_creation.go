package rx

import "sync"

// Merge creates a new channel that will merge the results of the given sources.
//
// If any of the sources emits an error, the returned channel will be immediately emits the error and close.
func Merge[T any](sources []<-chan Result[T], options ...Option) <-chan Result[T] {
	results := resultCh[T](options...)
	tmp := resultCh[T](options...)

	wg := new(sync.WaitGroup)
	wg.Add(len(sources))

	process := func(source <-chan Result[T]) {
		defer wg.Done()

		for v := range source {
			tmp <- v

			if v.IsError() {
				return
			}
		}
	}

	for _, source := range sources {
		go process(source)
	}

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

// Concat creates a new channel that will concatenate the results of the given sources
func Concat[T any](sources []<-chan Result[T], options ...Option) <-chan Result[T] {
	results := resultCh[T](options...)

	go func() {
		defer close(results)

		for _, source := range sources {
			for v := range source {
				results <- v

				if v.IsError() {
					return
				}
			}
		}
	}()

	return results
}
