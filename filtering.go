package rx

// TakeWhile takes values from the source channel while the predicate function returns false
func TakeWhile[T any](source <-chan Result[T], pred func(value T, index int) (bool, error), options ...Option) <-chan Result[T] {
	results := resultCh[T](prepend(WithBufferSize(cap(source)), options)...)

	go func() {
		defer close(results)

		index := 0
		for v := range source {
			v, err := v.Get()
			if err != nil {
				results <- Err[T](err)

				return
			}

			ok, err := pred(v, index)
			if err != nil {
				results <- Err[T](err)

				return
			}

			if !ok {
				return
			}

			results <- Ok(v)
			index++
		}
	}()

	return results
}

// Take takes the first n values from the source channel and emits them to the output rx.
// If the source channel emits less than n values, the output channel will close after emitting all values
func Take[T any](source <-chan Result[T], n int, options ...Option) <-chan Result[T] {
	return TakeWhile[T](source, func(_ T, index int) (bool, error) {
		return index < n, nil
	}, options...)
}

// Filter emits values from the source channel that pass the predicate function
func Filter[T any](source <-chan Result[T], predicate func(value T, index int) (bool, error), options ...Option) <-chan Result[T] {
	results := resultCh[T](prepend(WithBufferSize(cap(source)), options)...)

	go func() {
		defer close(results)

		index := 0
		for v := range source {
			v, err := v.Get()
			if err != nil {
				results <- Err[T](err)

				return
			}

			ok, err := predicate(v, index)
			if err != nil {
				results <- Err[T](err)

				return
			}

			if ok {
				results <- Ok(v)
			}

			index++
		}
	}()

	return results
}
