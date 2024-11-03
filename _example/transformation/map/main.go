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
	example4()
}

func example1() {
	println("Example 1")

	source := rx.Map(rx.Take(rx.Interval(500*time.Millisecond), 5, rx.WithBufferSize(2)), func(v int, _ int) (string, error) {
		return fmt.Sprintf("%04d", v+1), nil
	})

	println("capacity: ", cap(source))

	rx.Observe(source, rx.Observer[string]{
		Next: func(v string) {
			println("value: ", v)
		},
		Done: func() {
			println("example 1 done")
		},
	})
}

func example2() {
	println("Example 2: delay")

	source := rx.Map(rx.Take(rx.Interval(500*time.Millisecond), 5, rx.WithBufferSize(2)), func(v int, _ int) (string, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)

		return fmt.Sprintf("%04d", v+1), nil
	}, rx.WithBufferSize(0))

	println("capacity: ", cap(source))

	rx.Observe(source, rx.Observer[string]{
		Next: func(v string) {
			println("value: ", v)
		},
	})
}

func example3() {
	println("Example 3: error handling")

	source := rx.Map(rx.Take(rx.Interval(500*time.Millisecond), 5), func(v int, _ int) (string, error) {
		// simulate a slow operation
		<-time.After(1 * time.Second)

		if v == 3 {
			return "", fmt.Errorf("process error at v = %d", v)
		}

		return fmt.Sprintf("%04d", v+1), nil
	})

	rx.Observe(source, rx.Observer[string]{
		Next: func(v string) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
	})
}

func example4() {
	println("Example 4: error emitted from source channel")

	source := rx.Observable(func(observer rx.Observer[int]) {
		for i := 0; i < 5; i++ {
			observer.Next(i)

			<-time.After(1000 * time.Millisecond)
		}

		observer.Err(fmt.Errorf("source error"))
	})

	source2 := rx.Map(source, func(v int, _ int) (string, error) {
		return fmt.Sprintf("%04d", v+1), nil
	})

	rx.Observe(source2, rx.Observer[string]{
		Next: func(v string) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
		Done: func() {
			println("example 4 done")
		},
	})
}
