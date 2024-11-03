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
	example3()
}

func example1() {
	println("example 1")

	source := rx.Filter(rx.Take(rx.Interval(500*time.Millisecond), 5, rx.WithBufferSize(2)), func(v int, _ int) (bool, error) {
		return v%2 == 0, nil
	})

	println("capacity:", cap(source))

	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
	})
}

func example2() {
	println("Example 2: delay")

	source := rx.Filter(rx.Take(rx.Interval(500*time.Millisecond), 5, rx.WithBufferSize(2)), func(v int, _ int) (bool, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)

		return v%2 == 0, nil
	}, rx.WithBufferSize(1))

	println("capacity:", cap(source))

	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
	})
}

func example3() {
	println("Example 3")

	source := rx.Filter(rx.Take(rx.Interval(500*time.Millisecond), 5, rx.WithBufferSize(2)), func(v int, _ int) (bool, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)

		if v == 3 {
			return false, fmt.Errorf("process error at v = %d", v)
		}

		return v%2 == 0, nil
	}, rx.WithBufferSize(1))

	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
	})
}
