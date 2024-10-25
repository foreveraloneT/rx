// Package main is the entry point of the program
package main

import (
	"fmt"
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	println("example 1")
	ch1 := newSourceCh()
	out1, _ := rx.Filter(ch1, func(v int, _ int) (bool, error) {
		return v%2 == 0, nil
	})

	for v := range out1 {
		println("value: ", v)
	}

	println("example 2")
	ch2 := newSourceCh()
	out2, _ := rx.Filter(ch2, func(v int, _ int) (bool, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)

		return v%2 == 0, nil
	})
	for v := range out2 {
		println("value: ", v)
	}

	println("example 3: error handling")
	ch3 := newSourceCh()
	out3, errs := rx.Filter(ch3, func(v int, _ int) (bool, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)

		if v == 3 {
			return false, fmt.Errorf("process error")
		}

		return v%2 == 0, nil
	})

LOOP:
	for {
		select {
		case v, ok := <-out3:
			if !ok {
				break LOOP
			}
			println("value: ", v)
		case err := <-errs:
			if err != nil {
				println("error: ", err.Error())
			}
		}
	}
}

func newSourceCh() <-chan int {
	return rx.Take(rx.Interval(500*time.Millisecond), 5, rx.WithBufferSize(5))
}
