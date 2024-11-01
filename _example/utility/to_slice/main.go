// Package main is the entry point of the program
package main

import (
	"fmt"
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
}

func example1() {
	println("Example 1")

	ch := rx.Take(rx.Interval(400*time.Millisecond), 5)

	println("Consuming result. Please wait...")
	nn := rx.ToSlice(ch)

	fmt.Printf("Slice: %+v\n", nn)
}
