// Package main is the entry point of the program
package main

import (
	"fmt"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
	example2()
}

func example1() {
	fmt.Println("Example 1")

	source := rx.Generate(
		0,
		func(i int) bool { return i < 5 },
		func(i int) int { return i + 1 },
		func(i int) string { return fmt.Sprintf("%06d", i*1000) },
	)

	rx.Observe(source, rx.Observer[string]{
		Next: func(v string) {
			println("value:", v)
		},
	})
}

func example2() {
	fmt.Println("Example 2")

	source := rx.Generate(
		1,
		func(i int) bool { return i <= 11 },
		func(i int) int { return i + 2 },
		func(i int) int { return i },
		rx.WithBufferSize(1),
	)

	rx.Observe(source, rx.Observer[int]{
		Next: func(v int) {
			fmt.Println("value: ", v)
		},
		Done: func() {
			fmt.Println("done !")
		},
	})
}
