// Package main is the entry point of the program
package main

import (
	"errors"
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
	example2()
	example3()
}

func example1() {
	println("Example 1")

	ch := rx.Range(1, 10)
	out := rx.Reduce(ch, func(acc int, cur int, _ int) (int, error) {
		return acc + cur, nil
	}, 0)

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
		Done: func() {
			println("example 1 done")
		},
	})
}

func example2() {
	println("Example 2: delay")

	ch := rx.Range(1, 5)
	out := rx.Reduce(ch, func(acc int, cur int, _ int) (int, error) {
		// simulate a slow operation
		<-time.After(800 * time.Millisecond)

		return acc + cur, nil
	}, 0)

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
		Done: func() {
			println("example 2 done")
		},
	})
}

func example3() {
	println("Example 3: error handling")

	ch := rx.Range(1, 5)
	out := rx.Reduce(ch, func(acc int, cur int, index int) (int, error) {
		if index == 2 {
			return 0, errors.New("process error")
		}

		return acc + cur, nil
	}, 0)

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
		Done: func() {
			println("example 3 done")
		},
	})
}
