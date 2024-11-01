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
	out1, _ := rx.Observable[string](func(observer rx.Observer[string]) {
		defer observer.Done()

		observer.Next("Hello")
		<-time.After(500 * time.Millisecond)

		observer.Next("world")
		<-time.After(500 * time.Millisecond)

		observer.Next("my")
		<-time.After(500 * time.Millisecond)

		observer.Next("name")
		<-time.After(500 * time.Millisecond)

		observer.Next("is")
		<-time.After(500 * time.Millisecond)

		observer.Next("Kala")
		<-time.After(500 * time.Millisecond)
	}, rx.WithBufferSize(1))

	for v := range out1 {
		fmt.Println("value:", v)
	}
}

func example2() {
	println("Example 2: with error")

	out2, errs := rx.Observable[int](func(observer rx.Observer[int]) {
		defer observer.Done()

		observer.Next(1)
		<-time.After(500 * time.Millisecond)

		observer.Next(2)
		<-time.After(500 * time.Millisecond)

		observer.Err(fmt.Errorf("error appeared #1"))

		// These should not effect because the error has been emitted
		observer.Next(3)
		observer.Err(fmt.Errorf("error appeared #2"))
	})

	if err := rx.Observe(rx.Observer[int]{
		Next: func(v int) {
			fmt.Println("value:", v)
		},
		Err: func(err error) {
			fmt.Println("error:", err)
		},
	}, out2, errs); err != nil {
		fmt.Println("error from observation:", err)
	}
}
