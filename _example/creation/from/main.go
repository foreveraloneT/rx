// Package main is the entry point of the program
package main

import (
	"strings"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
}

func example1() {
	println("Example 1")

	source := rx.From(
		strings.Split("Hello, world. My name is Kala", " "),
		rx.WithBufferSize(2),
	)

	rx.Observe(source, rx.Observer[string]{
		Next: func(v string) {
			println("value:", v)
		},
	})
}
