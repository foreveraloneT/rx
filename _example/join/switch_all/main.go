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

	source := rx.Observable(func(observer rx.Observer[<-chan rx.Result[int]]) {
		defer observer.Done()

		source1 := rx.Interval(1000 * time.Millisecond)
		observer.Next(source1)

		<-time.After(3500 * time.Millisecond)

		source2 := rx.Map(rx.Interval(1000*time.Millisecond), func(v int, _ int) (int, error) {
			return v * 10, nil
		})
		observer.Next(source2)

		<-time.After(3000 * time.Millisecond)
		source3 := rx.Map(rx.Interval(2000*time.Millisecond), func(v int, _ int) (int, error) {
			return v * 100, nil
		})
		observer.Next(source3)
	})

	out := rx.Take(rx.SwitchAll(source), 30)

	// the source will be switched at 3.5s and 6.5s
	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
		Done: func() {
			println("example 1 done")
		},
	})
}
