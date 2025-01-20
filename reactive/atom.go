package reactive

import (
	"sync"

	"github.com/Olian04/reactive-go/reactive/internal"
	"github.com/google/uuid"
)

type Atom[T any] struct {
	id          string
	isDirty     bool
	value       T
	subscribers map[string]func()
	m           sync.Mutex
}

func NewAtom[T any](value T) *Atom[T] {
	return &Atom[T]{
		id:          uuid.New().String(),
		value:       value,
		isDirty:     true,
		subscribers: make(map[string]func()),
		m:           sync.Mutex{},
	}
}

func (a *Atom[T]) Set(value T) {
	a.m.Lock()
	defer a.m.Unlock()
	a.value = value
	a.isDirty = true
	for _, fn := range a.subscribers {
		fn()
	}
}

func (a *Atom[T]) Get() T {
	a.m.Lock()
	defer a.m.Unlock()
	id, markDirty := internal.RegisterAsDependency(func(id string) {
		a.subscribers[id] = nil
	})
	if id != "" {
		a.subscribers[id] = markDirty
	}
	return a.value
}
