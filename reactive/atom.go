package reactive

import (
	"github.com/Olian04/reactive-go/reactive/internal"
	"github.com/google/uuid"
)

type Atom[T any] struct {
	id          string
	isDirty     bool
	value       T
	subscribers map[string]func()
}

func NewAtom[T any](value T) *Atom[T] {
	return &Atom[T]{
		id:          uuid.New().String(),
		value:       value,
		isDirty:     true,
		subscribers: make(map[string]func()),
	}
}

func (a *Atom[T]) Set(value T) {
	a.value = value
	a.isDirty = true
	for _, fn := range a.subscribers {
		fn()
	}
}

func (a *Atom[T]) Get() T {
	id, markDirty := internal.AddDependency(a.id, func() bool {
		return a.isDirty
	})
	if id != "" {
		a.subscribers[id] = markDirty
	}
	return a.value
}
