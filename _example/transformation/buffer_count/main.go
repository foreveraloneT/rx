// Package main is the entry point of the program
package main

import (
	"fmt"
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
	example2()
}

func example1() {
	println("Example 1")

	source := rx.Take(rx.Interval(200*time.Millisecond), 15)
	buffered := rx.BufferCount(source, 4)

	rx.Observe(buffered, rx.Observer[[]int]{
		Next: func(v []int) {
			fmt.Printf("value: %+v\n", v)
		},
	})
}

func example2() {
	println("Example 2: error emitted from source channel")

	// TODO: Implement this example
}
