// Package main is the entry point of the program
package main

import (
	"strings"

	"github.com/foreveraloneT/rx"
)

func main() {
	ch := rx.From(
		strings.Split("Hello, world. My name is Kala", " "),
		rx.WithBufferSize(2),
	)

	for v := range ch {
		println("value: ", v)
	}
}
