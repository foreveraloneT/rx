package rx

// Tap used to perform side effects for each value emitted by the source rx.
// The output channel will emit the same values as the source rx.
// Please note that any error that occurs in the observer function will not be handled
func Tap[T any](c <-chan T, observer func(value T, index int)) <-chan T {
	out := make(chan T, cap(c))

	go func() {
		defer close(out)

		index := 0
		for v := range c {
			currentValue := v
			currentIndex := index
			go observer(currentValue, currentIndex)

			out <- v

			index++
		}
	}()

	return out
}

// ToSlice converts the values from the source channel to a slice
func ToSlice[T any](c <-chan T) []T {
	var result []T

	for v := range c {
		result = append(result, v)
	}

	return result
}
