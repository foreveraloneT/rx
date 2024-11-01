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

	ch := rx.Take(rx.Interval(1*time.Second), 5)

	for v := range ch {
		println("value: ", v)
	}
}
