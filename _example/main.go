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
	source := rx.Observable[string](func(observer rx.Observer[string]) {
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

	for result := range source {
		// Not recommend, but if you sure that there is no error, you can use `Value` method
		fmt.Println("value:", result.Value())
	}
}

func example2() {
	println("Example 2: with error")

	source := rx.Observable[int](func(observer rx.Observer[int]) {
		defer observer.Done()

		observer.Next(1)
		<-time.After(1000 * time.Millisecond)

		observer.Next(2)
		<-time.After(1000 * time.Millisecond)

		observer.Err(fmt.Errorf("error appeared #1"))

		// These should not effect because the error has been emitted
		observer.Next(3)
		observer.Err(fmt.Errorf("error appeared #2"))
	})

	if err := rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			fmt.Println("value:", v)
		},
		Err: func(err error) {
			fmt.Println("error:", err)
		},
	}); err != nil {
		fmt.Println("error from observation:", err)
	}
}

func example3() {
	println("Example 3: Another way to deals with results")

	source := rx.Observable[int](func(observer rx.Observer[int]) {
		defer observer.Done()

		observer.Next(1)
		<-time.After(1000 * time.Millisecond)

		observer.Next(2)
		<-time.After(1000 * time.Millisecond)

		observer.Err(fmt.Errorf("error appeared #1"))

		// These should not effect because the error has been emitted
		observer.Next(3)
		observer.Err(fmt.Errorf("error appeared #2"))
	})

	for result := range source {
		v, err := result.Get()
		if err != nil {
			fmt.Println("error:", err.Error())

			return
		}

		fmt.Println("value:", v)
	}
}
