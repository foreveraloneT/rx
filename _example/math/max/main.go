// Package main is the entry point of the program
package main

import (
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

	source := rx.Of(10, 11, 2, -99, 9, -1, 0)
	out := rx.Max(source)

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
	})
}

func example2() {
	println("Example 2")

	source := rx.Of(1.1, 1.09, -1.09, 0.1)
	out := rx.Max(source)

	rx.Observe(out, rx.Observer[float64]{
		Next: func(v float64) {
			println("value: ", v)
		},
	})
}

func example3() {
	println("Example 3")

	source := rx.Of("antelope", "ant", "bat", "battle", "cup", "car")
	out := rx.Max(source)

	rx.Observe(out, rx.Observer[string]{
		Next: func(v string) {
			println("value: ", v)
		},
	})
}

func example4() {
	println("Example 4: source channel that not emit any value")

	source := rx.Empty[int]()
	out := rx.Max(source)

	rx.Observe(out, rx.Observer[int]{
		Next: func(v int) {
			println("value: ", v)
		},
		Err: func(err error) {
			println("error: ", err.Error())
		},
	})
}
