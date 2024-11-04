package rx

// Reduce reduces the values from the source channel into a single value that emit when the source channel closed using the reducer function
func Reduce[T any, R any](source <-chan Result[T], reducer func(acc R, cur T, index int) (R, error), seed R, options ...Option) <-chan Result[R] {
	return TakeLast(Scan(source, reducer, seed), 1, options...)
}
