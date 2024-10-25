// Package main is the entry point of the program
package main

import (
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	ch := rx.Take(rx.Interval(1*time.Second), 5)

	for v := range ch {
		println("value: ", v)
	}
}
