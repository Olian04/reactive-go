package selector

type Selector[T any] struct {
	ID      string
	isDirty bool
	value   T
	getter  func() T
}

func New[T any](id string, getter func() T) *Selector[T] {
	return &Selector[T]{ID: id, getter: getter, isDirty: true}
}

func (s *Selector[T]) Get() T {
	if s.isDirty {
		s.value = s.getter()
		s.isDirty = false
	}
	return s.value
}
