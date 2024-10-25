// Package main is the entry point of the program
package main

import (
	"github.com/foreveraloneT/rx"
)

func main() {
	ch := rx.Range(1, 10)

	for v := range ch {
		println("value: ", v)
	}
}
