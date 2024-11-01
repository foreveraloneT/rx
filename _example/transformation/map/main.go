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
	println("Example 1")

	ch1 := newSourceCh()
	out1, _ := rx.Map(ch1, func(v int, _ int) (string, error) {
		return fmt.Sprintf("%04d", v+1), nil
	})

	for v := range out1 {
		println("value: ", v)
	}
}

func example2() {
	println("Example 2: delay")

	ch2 := newSourceCh()
	out2, _ := rx.Map(ch2, func(v int, _ int) (string, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)
		return fmt.Sprintf("%04d", v+1), nil
	})

	for v := range out2 {
		println("value: ", v)
	}
}

func example3() {
	println("Example 3: error handling")

	ch3 := newSourceCh()
	out3, errs := rx.Map(ch3, func(v int, _ int) (string, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)

		if v == 3 {
			return "", fmt.Errorf("process error")
		}

		return fmt.Sprintf("%04d", v+1), nil
	})

	rx.Observe(rx.Observer[string]{
		Next: func(v string) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
	}, out3, errs)
}

func newSourceCh() <-chan int {
	return rx.Take(rx.Interval(500*time.Millisecond), 5, rx.WithBufferSize(5))
}
