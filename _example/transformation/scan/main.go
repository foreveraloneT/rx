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

	ch1 := rx.Range(2, 10)
	out1, _ := rx.Scan(ch1, func(acc int, cur int, _ int) (int, error) {
		return acc + cur, nil
	}, 0)

	for v := range out1 {
		println("value: ", v)
	}
}

func example2() {
	println("Example 2: delay")

	ch2 := rx.Range(2, 5)
	out2, _ := rx.Scan(ch2, func(acc int, cur int, _ int) (int, error) {
		// simulate a slow operation
		<-time.After(800 * time.Millisecond)

		return acc + cur, nil
	}, 0)

	for v := range out2 {
		println("value: ", v)
	}
}

func example3() {
	println("Example 3: error handling")

	ch3 := rx.Range(2, 5)
	out3, errs := rx.Scan(ch3, func(acc int, cur int, index int) (int, error) {
		if index == 2 {
			return 0, errors.New("process error")
		}

		return acc + cur, nil
	}, 0)

	rx.Observe(rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
		Done: func() {
			println("Observation done")
		},
	}, out3, errs)
}
