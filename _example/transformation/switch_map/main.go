// Package main is the entry point of the program
package main

import (
	"errors"
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	println("example 1")
	ch1 := rx.Take(rx.Interval(1*time.Second), 10)
	out1, _ := rx.SwitchMap(ch1, func(v int, _ int) (<-chan int, <-chan error) {
		return rx.Map(rx.From([]int{1, 10, 100}), func(w int, _ int) (int, error) {
			<-time.After(400 * time.Millisecond)
			return w * (v + 1), nil
		})

	})

	rx.Observe(out1, nil, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Done: func() {
			println("Example 1 done")
		},
	})

	println("example 2: error handler")
	ch2 := rx.Take(rx.Interval(1*time.Second), 10)
	out2, errs := rx.SwitchMap(ch2, func(v int, _ int) (<-chan int, <-chan error) {
		return rx.Map(rx.From([]int{1, 10, 100}), func(w int, _ int) (int, error) {
			<-time.After(400 * time.Millisecond)
			if v == 4 && w == 10 {
				return 0, errors.New("process error at v = 4, w = 10")
			}

			return w * (v + 1), nil
		})

	})

	rx.Observe(out2, errs, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
	})
}
