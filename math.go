package rx

import (
	"golang.org/x/exp/constraints"
)

// Reduce reduces the values from the source channel into a single value that emit when the source channel closed using the reducer function
func Reduce[T any, R any](c <-chan T, reducer func(acc R, cur T, index int) (R, error), seed R, options ...Option) (<-chan R, <-chan error) {
	opts := parseOption(options...)
	out := make(chan R, opts.bufferSize)
	errs := make(chan error)

	go func() {
		defer close(out)
		defer close(errs)

		acc := seed
		index := 0
		for v := range c {
			var err error
			acc, err = reducer(acc, v, index)
			if err != nil {
				errs <- err

				return
			}

			index++
		}

		out <- acc
	}()

	return out, errs
}

// Min emits the minimum value from the source channel when the source channel closed
func Min[T constraints.Ordered](c <-chan T, options ...Option) <-chan T {
	seed := <-c
	out, _ := Reduce(c, func(acc T, cur T, _ int) (T, error) {
		if cur < acc {
			return cur, nil
		}

		return acc, nil
	}, seed, options...)

	return out
}

// Max emits the maximum value from the source channel when the source channel closed
func Max[T constraints.Ordered](c <-chan T, options ...Option) <-chan T {
	seed := <-c
	out, _ := Reduce(c, func(acc T, cur T, _ int) (T, error) {
		if cur > acc {
			return cur, nil
		}

		return acc, nil
	}, seed, options...)

	return out
}

// CountBy emits the number of values from the source channel that satisfy the predicate function when the source channel closed
func CountBy[T any](pred func(value T, index int) (bool, error), c <-chan T, options ...Option) (<-chan int, <-chan error) {
	return Reduce(c, func(acc int, cur T, index int) (int, error) {
		ok, err := pred(cur, index)
		if err != nil {
			return 0, err
		}

		if ok {
			return acc + 1, nil
		}

		return acc, nil
	}, 0, options...)
}

// Count emits the number of values from the source channel when the source channel closed
func Count[T any](c <-chan T, options ...Option) <-chan int {
	out, _ := CountBy(func(_ T, _ int) (bool, error) {
		return true, nil
	}, c, options...)

	return out
}
