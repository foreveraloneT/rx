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
	out1, _ := rx.Map(ch1, func(v int, _ int) (string, error) {
		return fmt.Sprintf("%04d", v+1), nil
	})

	for v := range out1 {
		println("value: ", v)
	}

	println("example 2")
	ch2 := newSourceCh()
	out2, _ := rx.Map(ch2, func(v int, _ int) (string, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)
		return fmt.Sprintf("%04d", v+1), nil
	})

	for v := range out2 {
		println("value: ", v)
	}

	println("example 3: error handling")
	ch3 := newSourceCh()
	out3, errs := rx.Map(ch3, func(v int, _ int) (string, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)

		if v == 3 {
			return "", fmt.Errorf("process error")
		}

		return fmt.Sprintf("%04d", v+1), nil
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
