package rx

// Reduce reduces the values from the source channel into a single value that emit when the source channel closed using the reducer function
func Reduce[T any, R any](source <-chan Result[T], reducer func(acc R, cur T, index int) (R, error), seed R, options ...Option) <-chan Result[R] {
	results := resultCh[R](options...)

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

			index++
		}

		results <- Ok(acc)
	}()

	return results
}
