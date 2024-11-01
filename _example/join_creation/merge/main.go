// Package main is the entry point of the program
package main

import (
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	example1()
	example2()
}

func example1() {
	println("Example 1: asynchronous")

	ch11 := rx.Take(rx.Interval(1*time.Second), 5)
	ch12, _ := rx.Map(
		rx.Take(rx.Interval(2*time.Second), 3),
		func(v int, _ int) (int, error) {
			return (v + 1) * 1000, nil
		},
	)

	out1 := rx.Merge([]<-chan int{ch11, ch12})

	for v := range out1 {
		println("value: ", v)
	}

}

func example2() {
	println("Example 2: synchronous, it will generate in random order")

	ch21 := rx.From([]int{1, 2, 3, 4, 5})
	ch22 := rx.From([]int{1000, 2000, 3000})
	out2 := rx.Merge([]<-chan int{ch21, ch22})

	for v := range out2 {
		println("value: ", v)
	}
}
