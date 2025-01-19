package reactive

import (
	"github.com/google/uuid"

	"github.com/Olian04/reactive-go/reactive/internal"
)

type Selector[T any] struct {
	id           string
	isDirty      bool
	value        T
	getter       func() T
	dependencies map[string]func() bool
	subscribers  map[string]func()
}

func NewSelector[T any](getter func() T) *Selector[T] {
	return &Selector[T]{
		id:           uuid.New().String(),
		getter:       getter,
		isDirty:      true,
		dependencies: make(map[string]func() bool),
		subscribers:  make(map[string]func()),
	}
}

func (s *Selector[T]) setDirty() {
	s.isDirty = true
	for _, fn := range s.subscribers {
		fn()
	}
}

func (s *Selector[T]) Get() T {
	if s.isDirty {
		internal.PushExecutionStack(&internal.ExecutionFrame{
			AddDependency: func(id string, fn func() bool) (string, func()) {
				s.dependencies[id] = fn
				return s.id, func() {
					s.setDirty()
				}
			},
		})
		s.value = s.getter()
		internal.PopExecutionStack()
		s.isDirty = false
	}
	id, markDirty := internal.AddDependency(s.id, func() bool {
		return s.isDirty
	})
	if id != "" {
		s.subscribers[id] = markDirty
	}
	return s.value
}

func (s *Selector[T]) Subscribe(id string, fn func()) {
	s.subscribers[id] = fn
}

func (s *Selector[T]) Unsubscribe(id string) {
	delete(s.subscribers, id)
}
