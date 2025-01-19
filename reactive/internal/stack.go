package internal

import (
	"container/list"
	"fmt"
)

type Stack[T any] struct {
	stack *list.List
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		stack: list.New(),
	}
}

func (s *Stack[T]) Push(value T) {
	s.stack.PushBack(value)
}

func (s *Stack[T]) Pop() (T, error) {
	if s.stack.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	e := s.stack.Back()
	s.stack.Remove(e)
	return e.Value.(T), nil
}

func (s *Stack[T]) Peek() (T, error) {
	if s.stack.Len() == 0 {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	return s.stack.Back().Value.(T), nil
}
