package rx

// BufferCount buffers the source channel values until the buffer size is reached, then emits the buffer and starts a new buffer
func BufferCount[T any](source <-chan Result[T], n int, option ...Option) <-chan Result[[]T] {
	results := resultCh[[]T](prepend(WithBufferSize(cap(source)), option)...)

	go func() {
		defer close(results)

		buffer := make([]T, 0, n)
		for v := range source {
			value, err := v.Get()
			if err != nil {
				results <- Err[[]T](err)

				return
			}

			buffer = append(buffer, value)
			if len(buffer) == n {
				results <- Ok(buffer)

				buffer = make([]T, 0, n)
			}
		}

		if len(buffer) > 0 {
			results <- Ok(buffer)
		}
	}()

	return results
}

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

// FlatMap transforms the values from the source channel to another channel using the provided mapper and joins them using the provided join function
func FlatMap[T any, R any](
	source <-chan Result[T],
	mapper func(value T, index int) <-chan Result[R],
	joinFn func(sources <-chan Result[<-chan Result[R]], options ...Option) <-chan Result[R],
	options ...Option,
) <-chan Result[R] {
	finalOptions := prepend(WithBufferSize(cap(source)), options)
	sources := resultCh[<-chan Result[R]](finalOptions...)

	go func() {
		defer close(sources)

		index := 0
		for v := range source {
			value, err := v.Get()
			if err != nil {
				sources <- Err[<-chan Result[R]](err)

				return
			}

			out := mapper(value, index)
			sources <- Ok(out)

			index++
		}
	}()

	return joinFn(sources, finalOptions...)
}

// MergeMap transforms the values from the source channel to another channel using the provided mapper and `MergeAll` them
func MergeMap[T any, R any](source <-chan Result[T], mapper func(value T, index int) <-chan Result[R], options ...Option) <-chan Result[R] {
	return FlatMap[T, R](source, mapper, MergeAll, options...)
}

// SwitchMap transforms the values from the source channel to another channel using the provided function and `SwitchAll` them
func SwitchMap[T any, R any](source <-chan Result[T], mapper func(value T, index int) <-chan Result[R], options ...Option) <-chan Result[R] {
	return FlatMap[T, R](source, mapper, SwitchAll, options...)
}

// GroupedResults represents a return value of `GroupBy` function
type GroupedResults[T any, K comparable] struct {
	Results <-chan Result[T]
	Key     K
}

// GroupBy groups the values from the source channel based on the keySelector function
func GroupBy[T any, K comparable](source <-chan Result[T], keySelector func(v T, index int) (K, error), options ...Option) <-chan Result[GroupedResults[T, K]] {
	finalOptions := prepend(WithBufferSize(cap(source)), options)
	results := resultCh[GroupedResults[T, K]](finalOptions...)

	go func() {
		defer close(results)

		groups := make(map[K]chan Result[T])
		defer func(groups map[K]chan Result[T]) {
			for _, group := range groups {
				close(group)
			}
		}(groups)

		index := 0
		for v := range source {
			value, err := v.Get()
			if err != nil {
				results <- Err[GroupedResults[T, K]](err)

				return
			}

			key, err := keySelector(value, index)
			if err != nil {
				results <- Err[GroupedResults[T, K]](err)

				return
			}

			if _, ok := groups[key]; !ok {
				groups[key] = resultCh[T](finalOptions...)

				results <- Ok(GroupedResults[T, K]{
					Key:     key,
					Results: groups[key],
				})
			}

			groups[key] <- Ok(value)

			index++
		}
	}()

	return results
}

// Scan transforms the values from the source channel using the provided accumulator function (or reducer function).
// Like Reduce, but emits values every time the source channel emits a value.
func Scan[T any, R any](source <-chan Result[T], reducer func(acc R, cur T, index int) (R, error), seed R, options ...Option) <-chan Result[R] {
	results := resultCh[R](prepend(WithBufferSize(cap(source)), options)...)

	go func() {
		defer close(results)

		acc := seed
		index := 0
		for v := range source {
			value, err := v.Get()
			if err != nil {
				results <- Err[R](err)

				return
			}

			acc, err = reducer(acc, value, index)
			if err != nil {
				results <- Err[R](err)

				return
			}

			results <- Ok(acc)

			index++
		}
	}()

	return results
}
