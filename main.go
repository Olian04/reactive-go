package main

import (
	"fmt"

	"github.com/olian04/reactive-go/atom"
	"github.com/olian04/reactive-go/selector"
)

func main() {
	a := atom.New("a", 1)
	b := atom.New("b", 2)
	c := selector.New("c", func() int {
		return a.Get() + b.Get()
	})

	fmt.Println(c.Get())
}
