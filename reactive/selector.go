package reactive

import (
	"sync"

	"github.com/google/uuid"

	"github.com/Olian04/reactive-go/reactive/internal"
)

type Selector[T any] struct {
	id          string
	isDirty     bool
	value       T
	getter      func() T
	subscribers map[string]func()
	m           sync.Mutex
}

func NewSelector[T any](getter func() T) *Selector[T] {
	return &Selector[T]{
		id:          uuid.New().String(),
		getter:      getter,
		isDirty:     true,
		subscribers: make(map[string]func()),
		m:           sync.Mutex{},
	}
}

func (s *Selector[T]) Get() T {
	s.m.Lock()
	defer s.m.Unlock()
	if s.isDirty {
		internal.PushExecutionStack(&internal.ExecutionFrame{
			RegisterAsDependency: func(removeDependency func(string)) (string, func()) {
				return s.id, func() {
					s.isDirty = true
					for _, fn := range s.subscribers {
						fn()
					}
				}
			},
		})
		s.value = s.getter()
		internal.PopExecutionStack()
		s.isDirty = false
	}
	id, markDirty := internal.RegisterAsDependency(func(id string) {
		s.subscribers[id] = nil
	})
	if id != "" {
		s.subscribers[id] = markDirty
	}
	return s.value
}
