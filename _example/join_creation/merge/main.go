// Package main is the entry point of the program
package main

import (
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
	example2()
	example3()
}

func example1() {
	println("Example 1: asynchronous")

	ch1 := rx.Take(rx.Interval(1*time.Second), 5)
	ch2 := rx.Map(
		rx.Take(rx.Interval(2*time.Second), 3),
		func(v int, _ int) (int, error) {
			return (v + 1) * 1000, nil
		},
	)

	out := rx.Merge([]<-chan rx.Result[int]{ch1, ch2})

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
	})
}

func example2() {
	println("Example 2: synchronous, it will generate in random order")

	ch1 := rx.From([]int{1, 2, 3, 4, 5})
	ch2 := rx.From([]int{1000, 2000, 3000})
	out := rx.Merge([]<-chan rx.Result[int]{ch1, ch2})

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
	})
}

func example3() {
	println("Example 3: there are some errors")

	// TODO: add example 3
}
