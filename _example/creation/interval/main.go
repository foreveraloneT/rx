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

	source := rx.Interval(1 * time.Second)

	println("[Warning]: The program is never ending, press Ctrl+C to stop")
	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			println(v)
		},
		Err: func(err error) {
			println(err.Error())
		},
	})
}
