// Package main is the entry point of the program
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
	example2()
}

func example1() {
	println("Example 1")

	source := rx.ToSlice(rx.Of("Hello", "world", "my", "name", "is", "Kala"))

	rx.Observe(source, rx.Observer[[]string]{
		Next: func(v []string) {
			fmt.Println("value:", strings.Join(v, ", "))
		},
		Err: func(err error) {
			fmt.Println("error:", err)
		},
		Done: func() {
			fmt.Println("example 1 done")
		},
	})
}

func example2() {
	println("Example 2")

	source := rx.Observable(func(observer rx.Observer[string]) {
		defer observer.Done()

		observer.Next("Hello")
		<-time.After(500 * time.Millisecond)
		observer.Next("world")
		<-time.After(500 * time.Millisecond)
		observer.Err(fmt.Errorf("error appeared"))
	})

	rx.Observe(rx.ToSlice(source), rx.Observer[[]string]{
		Next: func(v []string) {
			fmt.Println("value:", strings.Join(v, ", "))
		},
		Err: func(err error) {
			fmt.Println("error:", err)
		},
		Done: func() {
			fmt.Println("example 2 done")
		},
	})
}
