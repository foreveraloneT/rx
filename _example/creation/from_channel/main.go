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
	println("Example 1, default buffer size equals to input channel")

	ch := make(chan int, 3)
	go func() {
		defer close(ch)

		for i := 0; i < 5; i++ {
			<-time.After(500 * time.Millisecond)

			ch <- i
		}
	}()

	source := rx.FromChannel(ch)
	println("buffer size:", cap(source))

	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			fmt.Println("value:", v)
		},
	})
}

func example2() {
	println("Example 2, custom buffer size")

	ch := make(chan int, 3)
	go func() {
		defer close(ch)

		for i := 0; i < 5; i++ {
			<-time.After(500 * time.Millisecond)

			ch <- i
		}
	}()

	source := rx.FromChannel(ch, rx.WithBufferSize(0))
	println("buffer size:", cap(source))

	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			fmt.Println("value:", v)
		},
	})
}
