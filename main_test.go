package main

import (
	"testing"

	"github.com/Olian04/reactive-go/reactive"
)

func TestMain(t *testing.T) {
	count := 0
	a := reactive.NewAtom(1)
	b := reactive.NewAtom(2)
	c := reactive.NewSelector(func() int {
		count++
		return a.Get() + b.Get()
	})

	if v := c.Get(); v != 3 {
		t.Fatalf("Expected 3, got %d", v)
	}

	a.Set(3)
	if v := c.Get(); v != 5 {
		t.Fatalf("Selector was not recomputed after atom a was set. Expected 5, got %d", v)
	}

	if count != 2 {
		t.Fatalf("Selector was not marked as dirty after atom a was set")
	}
}
