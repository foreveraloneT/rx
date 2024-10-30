package rx

// TakeWhile takes values from the source channel while the predicate function returns false
func TakeWhile[T any](c <-chan T, pred func(value T, index int) (bool, error), options ...Option) (<-chan T, <-chan error) {
	out, errs := observableChWithErrs[T](options...)

	go func() {
		defer close(out)
		defer close(errs)

		index := 0
		for v := range c {
			ok, err := pred(v, index)
			if err != nil {
				errs <- err

				return
			}

			if !ok {
				return
			}

			out <- v
			index++
		}
	}()

	return out, errs
}

// Take takes the first n values from the source channel and emits them to the output rx.
// If the source channel emits less than n values, the output channel will close after emitting all values
func Take[T any](c <-chan T, n int, options ...Option) <-chan T {
	out := observableCh[T](options...)

	go func() {
		defer close(out)

		count := 0
		for v := range c {
			if count == n {
				return
			}

			out <- v
			count++
		}
	}()

	return out
}

// Filter filters the values from the source channel based on the predicate function.
func Filter[T any](c <-chan T, pred func(value T, index int) (bool, error), options ...Option) (<-chan T, <-chan error) {
	out, errs := observableChWithErrs[T](options...)

	go func() {
		defer close(out)
		defer close(errs)

		index := 0
		for v := range c {
			ok, err := pred(v, index)
			if err != nil {
				errs <- err

				return
			}

			if ok {
				out <- v
			}

			index++
		}
	}()

	return out, errs
}
