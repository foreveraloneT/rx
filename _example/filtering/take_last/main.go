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

	source := rx.Map(rx.Range(1, 20), func(v int, _ int) (int, error) {
		<-time.After(200 * time.Millisecond)

		return v, nil
	})
	lastThree := rx.TakeLast(source, 3)

	rx.Observe(lastThree, rx.Observer[int]{
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
	println("Example 2: error emitted from source")

	source := rx.Map(rx.Range(1, 20), func(v int, _ int) (int, error) {
		<-time.After(200 * time.Millisecond)
		if v == 20 {
			return 0, fmt.Errorf("error at v = %d", v)
		}

		return v, nil
	})
	lastThree := rx.TakeLast(source, 3)

	rx.Observe(lastThree, rx.Observer[int]{
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
