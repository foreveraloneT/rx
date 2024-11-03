// Package main is the entry point of the program
package main

import (
	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
}

func example1() {
	println("Example 1")

	source := rx.Range(1, 10)

	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
	})
}
