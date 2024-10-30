// Package main is the entry point of the program
package main

import (
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	source := rx.From([]int{1, 3, 5, 7, 9, 10, 11})
	tmp, _ := rx.Map(source, func(v int, _ int) (int, error) {
		<-time.After(1000 * time.Millisecond)

		return v, nil
	})

	out, _ := rx.TakeWhile(tmp, func(v int, index int) (bool, error) {
		return v%2 != 0, nil
	})

	rx.Observe(out, nil, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Done: func() {
			println("channel closed")
		},
	})
}
