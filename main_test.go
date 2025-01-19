package main

import (
	"fmt"
	"testing"

	"github.com/Olian04/reactive-go/atom"
	"github.com/Olian04/reactive-go/selector"
)

func TestMain(t *testing.T) {
	a := atom.New("a", 1)
	b := atom.New("b", 2)
	c := selector.New("c", func() int {
		fmt.Println("selector c is dirty")
		return a.Get() + b.Get()
	})

	if v := c.Get(); v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}
	a.Set(3)
	if v := c.Get(); v != 5 {
		t.Fatalf("Selector was not marked as dirty after atom a was set")
	}
}
