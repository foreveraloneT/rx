package rx

// CatchError catches errors from the source channel.
// If errors are emitted, it will call the selector function to get a new channel to observe
func CatchError[T any](
	source <-chan Result[T],
	selector func(err error) <-chan Result[T],
	options ...Option,
) <-chan Result[T] {
	results := resultCh[T](prepend(WithBufferSize(cap(source)), options)...)

	go func() {
		defer close(results)

		var err error
		for v := range source {
			var value T
			value, err = v.Get()
			if err != nil {
				break
			}

			results <- Ok(value)
		}

		if err == nil {
			return
		}

		for v := range selector(err) {
			results <- v
		}
	}()

	return results
}
