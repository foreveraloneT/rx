package rx

// Map transforms the values from the source channel using the mapper function
func Map[T any, R any](source <-chan Result[T], mapper func(value T, index int) (R, error), options ...Option) <-chan Result[R] {
	results := resultCh[R](prepend(WithBufferSize(cap(source)), options)...)

	go func() {
		defer close(results)

		index := 0
		for v := range source {
			value, err := v.Get()
			if err != nil {
				results <- Err[R](err)

				return
			}

			mapped, err := mapper(value, index)
			if err != nil {
				results <- Err[R](err)

				return
			}

			results <- Ok(mapped)

			index++
		}
	}()

	return results
}
