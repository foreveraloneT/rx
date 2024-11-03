// Package main is the entry point of the program
package main

import (
	"fmt"
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
	example2()
}

func example1() {
	println("Example 1")

	ch := rx.Take(rx.Interval(1*time.Second), 10)
	out := rx.MergeMap(ch, func(v int, _ int) <-chan rx.Result[int] {
		return rx.Map(rx.From([]int{1, 10, 100}), func(w int, _ int) (int, error) {
			<-time.After(400 * time.Millisecond)
			return w * (v + 1), nil
		})
	})

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Done: func() {
			println("Example 1 done")
		},
	})
}

func example2() {
	println("Example 2: error handler")

	ch := rx.Take(rx.Interval(1*time.Second), 10)
	out := rx.MergeMap(ch, func(v int, _ int) <-chan rx.Result[int] {
		return rx.Map(rx.From([]int{1, 10, 100}), func(w int, _ int) (int, error) {
			<-time.After(400 * time.Millisecond)
			if v == 4 && w == 10 {
				return 0, fmt.Errorf("process error at v = %d, w = %d", v, w)
			}

			return w * (v + 1), nil
		})

	})

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
	})
}
