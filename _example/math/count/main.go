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
	println("Example 1: channel that closed without emitting any value")

	count := rx.Count(rx.Empty[int]())

	rx.Observe(count, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
		Err: func(err error) {
			println("error:", err.Error())
		},
	})
}

func example2() {
	println("Example 2")

	count := rx.Count(rx.Of[int](1, 2, 3, 4))

	rx.Observe(count, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
		Err: func(err error) {
			println("error:", err.Error())
		},
	})
}

func example3() {
	println("Example 3: error handler")

	source := rx.Observable(func(observer rx.Observer[int]) {
		observer.Next(1)
		<-time.After(1 * time.Second)

		observer.Next(2)
		<-time.After(1 * time.Second)

		observer.Next(3)
		<-time.After(1 * time.Second)

		observer.Err(errors.New("something went wrong"))
	})

	count := rx.Count(source)

	rx.Observe(count, rx.Observer[int]{
		Next: func(v int) {
			println("value:", v)
		},
		Err: func(err error) {
			println("error:", err.Error())
		},
	})
}
