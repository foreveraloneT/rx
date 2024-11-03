// Package main is the entry point of the program
package main

import "github.com/foreveraloneT/rx"

func main() {
	example1()
}

func example1() {
	println("Example 1")

	source := rx.Empty[int]()

	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
		Err: func(err error) {
			println("error:", err.Error())
		},
		Done: func() {
			println("example 1 done")
		},
	})
}
