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

	ch := rx.Take(rx.Interval(400*time.Millisecond), 10)

	tapped := rx.Tap(ch, func(value int, index int) {
		fmt.Printf("Tapping value: %v at index: %v\n", value, index)
	})

	mapped := rx.Map(tapped, func(value int, _ int) (int, error) {
		return value * 10, nil
	})

	rx.Observe(mapped, rx.Observer[int]{
		Next: func(v int) {
			fmt.Println("value: ", v)
		},
	})
}
