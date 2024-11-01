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

	ch1 := rx.Take(rx.Interval(1*time.Second), 10)
	out1, _ := rx.SwitchMap(ch1, func(v int, _ int) (<-chan int, <-chan error) {
		return rx.Map(rx.From([]int{1, 10, 100}), func(w int, _ int) (int, error) {
			<-time.After(400 * time.Millisecond)
			return w * (v + 1), nil
		})

	})

	rx.Observe(rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Done: func() {
			println("Example 1 done")
		},
	}, out1)
}

func example2() {
	println("Example 2: error handler")

	ch2 := rx.Take(rx.Interval(1*time.Second), 10)
	out2, errs := rx.SwitchMap(ch2, func(v int, _ int) (<-chan int, <-chan error) {
		return rx.Map(rx.From([]int{1, 10, 100}), func(w int, _ int) (int, error) {
			<-time.After(400 * time.Millisecond)
			if v == 4 && w == 10 {
				return 0, fmt.Errorf("process error at v = %d, w = %d", v, w)
			}

			return w * (v + 1), nil
		})

	})

	rx.Observe(rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
	}, out2, errs)
}
