// Package main is the entry point of the program
package main

import (
	"errors"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
}

func example1() {
	println("Example 1")

	source := rx.ThrowError[string](errors.New("this is an error"))

	rx.Observe(source, rx.Observer[string]{
		Next: func(v string) {
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
