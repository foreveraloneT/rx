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

	ch := rx.Take(rx.Interval(200*time.Millisecond), 15)
	out := rx.BufferCount(ch, 4)

	for v := range out {
		fmt.Printf("value: %v\n", v)
	}
}
