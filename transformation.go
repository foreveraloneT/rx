package rx

// BufferCount buffers the source channel values until the buffer size is reached, then emits the buffer and starts a new buffer
func BufferCount[T any](c <-chan T, n int, option ...Option) <-chan []T {
	out := observableCh[[]T](option...)

	go func() {
		defer close(out)

		buffer := make([]T, 0, n)
		for v := range c {
			buffer = append(buffer, v)
			if len(buffer) == n {
				out <- buffer
				buffer = make([]T, 0, n)
			}
		}

		if len(buffer) > 0 {
			out <- buffer
		}
	}()

	return out
}

// Map transforms the values from the source channel using the provided function
func Map[T any, R any](c <-chan T, iter func(value T, index int) (R, error), options ...Option) (<-chan R, <-chan error) {
	out, errs := observableChWithErrs[R](options...)

	go func() {
		defer close(out)
		defer close(errs)

		index := 0
		for v := range c {
			r, err := iter(v, index)
			if err != nil {
				errs <- err

				return
			}

			out <- r

			index++
		}
	}()

	return out, errs
}

// Scan transforms the values from the source channel using the provided accumulator function (or reducer function).
// Like Reduce, but emits values every time the source channel emits a value.
func Scan[T any, R any](c <-chan T, accumulator func(acc R, cur T, index int) (R, error), seed R, options ...Option) (<-chan R, <-chan error) {
	out, errs := observableChWithErrs[R](options...)

	go func() {
		defer close(out)
		defer close(errs)

		acc := seed
		index := 0
		for v := range c {
			var err error
			acc, err = accumulator(acc, v, index)
			if err != nil {
				errs <- err

				return
			}

			out <- acc

			index++
		}
	}()

	return out, errs
}

// MergeMap transforms the values from the source channel to another channel using the provided function and MergeAll them
func MergeMap[T any, R any](c <-chan T, iter func(value T, index int) (<-chan R, <-chan error), options ...Option) (<-chan R, <-chan error) {
	out, errs := observableChWithErrs[R](options...)

	outChs := make(chan (<-chan R))
	errsChs := make(chan (<-chan error))
	flattenOut := MergeAll(outChs, options...)
	flattenErrs := MergeAll(errsChs)

	go func() {
		defer close(outChs)
		defer close(errsChs)

		index := 0
		for v := range c {
			outCh, errsCh := iter(v, index)
			outChs <- outCh
			errsChs <- errsCh

			index++
		}
	}()

	go func() {
		defer close(out)
		defer close(errs)

		Observe(Observer[R]{
			Next: func(v R) {
				out <- v
			},
			Err: func(err error) {
				errs <- err
			},
		}, flattenOut, flattenErrs)
	}()

	return out, errs
}

// SwitchMap transforms the values from the source channel to another channel using the provided function and SwitchAll them
func SwitchMap[T any, R any](c <-chan T, iter func(value T, index int) (<-chan R, <-chan error), options ...Option) (<-chan R, <-chan error) {
	out, errs := observableChWithErrs[R](options...)

	outChs := make(chan (<-chan R))
	errsChs := make(chan (<-chan error))
	flattenOut := SwitchAll(outChs, options...)
	flattenErrs := SwitchAll(errsChs)

	go func() {
		defer close(outChs)
		defer close(errsChs)

		index := 0
		for v := range c {
			outCh, errsCh := iter(v, index)
			outChs <- outCh
			errsChs <- errsCh

			index++
		}
	}()

	go func() {
		defer close(out)
		defer close(errs)

		Observe(Observer[R]{
			Next: func(v R) {
				out <- v
			},
			Err: func(err error) {
				errs <- err
			},
		}, flattenOut, flattenErrs)
	}()

	return out, errs
}

// GroupByObserver represents a result of GroupBy
type GroupByObserver[T any, K any] struct {
	Val <-chan T
	Key K
}

// GroupBy groups the values from the source channel based on the keySelector function
func GroupBy[T any, K comparable](c <-chan T, keySelector func(v T, index int) (K, error), options ...Option) (<-chan GroupByObserver[T, K], chan error) {
	out, errs := observableChWithErrs[GroupByObserver[T, K]](options...)

	go func() {
		defer close(out)
		defer close(errs)

		groups := make(map[K]chan T)
		defer func(groups map[K]chan T) {
			for _, group := range groups {
				close(group)
			}
		}(groups)

		index := 0
		for v := range c {
			key, err := keySelector(v, index)
			if err != nil {
				errs <- err

				return
			}

			if _, ok := groups[key]; !ok {
				groups[key] = observableCh[T](options...)

				out <- GroupByObserver[T, K]{
					Val: groups[key],
					Key: key,
				}
			}

			groups[key] <- v

			index++
		}
	}()

	return out, errs
}
