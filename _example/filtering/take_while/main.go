// Package main is the entry point of the program
package main

import (
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
}

func example1() {
	println("Example 1")

	source1 := rx.Map(rx.From([]int{1, 3, 5, 7, 9, 10, 11, 12, 13}), func(v int, _ int) (int, error) {
		<-time.After(1000 * time.Millisecond)

		return v, nil
	})

	source2 := rx.TakeWhile(source1, func(v int, index int) (bool, error) {
		return v%2 != 0, nil
	})

	rx.Observe(source2, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Done: func() {
			println("channel closed")
		},
	})
}
