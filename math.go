package rx

import (
	"errors"

	"golang.org/x/exp/constraints"
)

// Reduce reduces the values from the source channel into a single value that emit when the source channel closed using the reducer function
// If the source channel is empty, it will lead ro an `ErrNotFound`
func Reduce[T any, R any](source <-chan Result[T], reducer func(acc R, cur T, index int) (R, error), seed R, options ...Option) <-chan Result[R] {
	pred := func(_ R, _ int) (bool, error) {
		return true, nil
	}

	return Last(Scan(source, reducer, seed), pred, options...)
}

// Min emits the minimum value from the source channel when the source channel closed
// If the source channel is empty, it will lead ro an `ErrNotFound`
func Min[T constraints.Ordered](source <-chan Result[T], options ...Option) <-chan Result[T] {
	return Reduce(source, func(acc T, cur T, index int) (T, error) {
		if index == 0 {
			return cur, nil
		}

		if cur < acc {
			return cur, nil
		}

		return acc, nil
	}, empty[T](), options...)
}

// Max emits the maximum value from the source channel when the source channel closed
// If the source channel is empty, it will lead ro an `ErrNotFound`
func Max[T constraints.Ordered](source <-chan Result[T], options ...Option) <-chan Result[T] {
	return Reduce(source, func(acc T, cur T, index int) (T, error) {
		if index == 0 {
			return cur, nil
		}

		if cur > acc {
			return cur, nil
		}

		return acc, nil
	}, empty[T](), options...)
}

// CountBy emits the number of values from the source channel that satisfy the predicate function when the source channel closed
func CountBy[T any](source <-chan Result[T], pred func(value T, index int) (bool, error), options ...Option) <-chan Result[int] {
	return Reduce(source, func(acc int, cur T, index int) (int, error) {
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
func Count[T any](source <-chan Result[T], options ...Option) <-chan Result[int] {
	count := CountBy(source, func(_ T, _ int) (bool, error) {
		return true, nil
	}, options...)

	return CatchError(count, func(err error) <-chan Result[int] {
		if errors.Is(err, ErrNotFound) {
			return Of(0)
		}

		return ThrowError[int](err)
	})
}
