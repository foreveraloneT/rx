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

	println("Wait for 6 seconds, then number 0 will be printed")
	source := rx.Timer(6 * time.Second)
	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			println(v)
		},
	})
}
