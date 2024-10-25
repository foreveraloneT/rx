// Package main is the entry point of the program
package main

import (
	"time"

	"github.com/foreveraloneT/rx"
)

func main() {
	ch := rx.Interval(1 * time.Second)

	println("[Warning]: The program is never ending, press Ctrl+C to stop")
	for v := range ch {
		println("value: ", v)
	}
}
