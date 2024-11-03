package rx

// Tap used to perform side effects for each value emitted by the source rx.
// The output channel will emit the same values as the source rx.
// Please note that any error that occurs in the `observer` function will not be handled
func Tap[T any](source <-chan Result[T], observer func(value T, index int), options ...Option) <-chan Result[T] {
	results := resultCh[T](prepend(WithBufferSize(cap(source)), options)...)

	go func() {
		defer close(results)

		index := 0
		for v := range source {
			currentValue := v
			currentIndex := index

			v, err := currentValue.Get()
			if err == nil {
				go observer(v, currentIndex)
			}

			results <- currentValue

			index++
		}
	}()

	return results
}

// ToSlice waits until the source channel is closed and then emits a single slice containing all the values emitted by the source channel.
// If the source channel emits an error, the output channel will not emit any value and will be closed immediately.
func ToSlice[T any](source <-chan Result[T], options ...Option) <-chan Result[[]T] {
	result := resultCh[[]T](options...)

	go func() {
		defer close(result)

		var values []T
		for v := range source {
			value, err := v.Get()
			if err != nil {
				return
			}

			values = append(values, value)
		}

		result <- Ok(values)
	}()

	return result
}
