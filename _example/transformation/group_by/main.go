// Package main is the entry point of the program
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/foreveraloneT/rx"
)

type programmingTool struct {
	Type string
	Name string
}

type group[T any] struct {
	Key    string
	Values []T
}

func main() {
	example1()
}

func example1() {
	println("Example 1")

	source := newSource()

	grouped := rx.GroupBy(source, func(v programmingTool, _ int) (string, error) {
		return v.Type, nil
	}, rx.WithBufferSize(1))

	out := rx.MergeMap(grouped, func(groupedResults rx.GroupedResults[programmingTool, string], _ int) <-chan rx.Result[*group[string]] {
		return rx.Reduce(groupedResults.Results, func(acc *group[string], v programmingTool, _ int) (*group[string], error) {
			if v.Name == "Vite" {
				<-time.After(2 * time.Second)

				return nil, fmt.Errorf("process error at name = %s", v.Name)
			}

			return &group[string]{
				Key:    acc.Key,
				Values: append(acc.Values, v.Name),
			}, nil
		}, &group[string]{
			Key:    groupedResults.Key,
			Values: make([]string, 0),
		})
	})

	if err := rx.Observe(out, rx.Observer[*group[string]]{
		Next: func(v *group[string]) {
			fmt.Printf("Type: %s, Items: [%s]\n", v.Key, strings.Join(v.Values, ", "))
		},
		Err: func(err error) {
			fmt.Println("Error:", err.Error())
		},
		Done: func() {
			fmt.Println("Done")
		},
	}); err != nil {
		fmt.Println("error from observation:", err.Error())
	}
}

func newSource() <-chan rx.Result[programmingTool] {
	pts := []programmingTool{
		{
			Type: "ide",
			Name: "VSCode",
		},
		{
			Type: "language",
			Name: "Go",
		},
		{
			Type: "framework",
			Name: "React",
		},
		{
			Type: "ide",
			Name: "Cursor",
		},
		{
			Type: "language",
			Name: "Typescript",
		},
		{
			Type: "bundler",
			Name: "Vite",
		},
		{
			Type: "framework",
			Name: "Vue",
		},
		{
			Type: "bundler",
			Name: "Webpack",
		},
	}

	return rx.Observable[programmingTool](func(observer rx.Observer[programmingTool]) {
		defer observer.Done()

		for _, pt := range pts {
			observer.Next(pt)

			<-time.After(400 * time.Millisecond)
		}
	}, rx.WithBufferSize(1))
}
