package main

import (
	"fmt"
	"testing"

	"github.com/Olian04/reactive-go/reactive"
)

func TestSelectorWithMultipleDependencies(t *testing.T) {
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

func TestSelectorWithSelectorDependency(t *testing.T) {
	a := reactive.NewAtom(1)
	b := reactive.NewSelector(func() int {
		return a.Get() + 1
	})
	c := reactive.NewSelector(func() int {
		return b.Get() + 1
	})

	if v := a.Get(); v != 1 {
		t.Fatalf("Expected 1, got %d", v)
	}

	a.Set(3)
	if v := c.Get(); v != 5 {
		t.Fatalf("Selector was not recomputed after atom a was set. Expected 5, got %d", v)
	}
}

func TestSelectorCachingRecomputation(t *testing.T) {
	count := 0
	a := reactive.NewAtom(1)
	b := reactive.NewSelector(func() int {
		count++
		return a.Get() + 1
	})
	c := reactive.NewSelector(func() int {
		return b.Get() + 1
	})
	d := reactive.NewSelector(func() int {
		return b.Get() + 1
	})

	b.Get()
	c.Get()
	d.Get()

	if count != 1 {
		t.Fatalf("Selector was recomputed more than once")
	}

	a.Set(2)
	d.Get()
	c.Get()
	b.Get()

	if count != 2 {
		t.Fatalf("Selector was recomputed more than once after updating dependency")
	}
}

func TestWithMultipleGoRoutines(t *testing.T) {
	testSize := 100 // TODO: Fix this race condition
	// Is OK for small test sizes (like 5)
	// But fails with multiple different errors for larger test sizes (like 100)

	ch := make(chan error, testSize)

	for range testSize {
		go func() {
			a := reactive.NewAtom(1)
			c := reactive.NewSelector(func() int {
				return a.Get()
			})

			if v := c.Get(); v != 1 {
				ch <- fmt.Errorf("Expected 1, got %d", v)
				return
			}

			a.Set(3)
			if v := c.Get(); v != 3 {
				ch <- fmt.Errorf("Selector was not recomputed after atom a was set. Expected 3, got %d", v)
				return
			}
			ch <- nil
		}()
	}

	count := 0
	for range testSize {
		if err := <-ch; err != nil {
			t.Fatal(err)
		}
		count++
	}

	if count != testSize {
		t.Fatalf("Expected %d goroutines to finish, got %d", testSize, count)
	}
}
