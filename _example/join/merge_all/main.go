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

	sources := rx.Observable(func(observer rx.Observer[<-chan rx.Result[int]]) {
		defer observer.Done()

		source1 := rx.Interval(1000 * time.Millisecond)
		observer.Next(source1) // source 1 will emit at 1000, 2000, 3000, ...

		<-time.After(1000 * time.Millisecond)

		source2 := rx.Map(rx.Interval(2000*time.Millisecond), func(v int, _ int) (int, error) {
			return v * 10, nil
		})
		observer.Next(source2) //  source 2 will emit at 3000, 5000, 7000, ...

		<-time.After(1000 * time.Millisecond)
		source3 := rx.Map(rx.Interval(1000*time.Millisecond), func(v int, _ int) (int, error) {
			return v * 100, nil
		})
		observer.Next(source3) // source  3 will emit at 3000, 4000, 5000, ...
	})

	/*
		Table to show the example of relationship between time and value emitting from each sources

		| At time (second) | source 1 | source 2 | source 3 |
		|------------------|----------|----------|----------|
		| 1                | 0        |          |          |
		| 2                | 1        |          |          |
		| 3                | 2        | 0        | 0        |
		| 4                | 3        | 100      |          |
		| 5                | 4        | 10       | 200      |
		| 6                | 5        |          | 300      |
		| 7                | 6        | 20       | 400      |
	*/

	out := rx.Take(rx.MergeAll(sources), 30)

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
		Done: func() {
			println("example 1 done")
		},
	})
}
