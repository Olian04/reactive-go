package reactive

type Atom[T any] struct {
	ID      string
	isDirty bool
	value   T
}

func New[T any](id string, value T) *Atom[T] {
	return &Atom[T]{ID: id, value: value, isDirty: true}
}

func (a *Atom[T]) Set(value T) {
	a.value = value
	a.isDirty = true
}

func (a *Atom[T]) Get() T {
	return a.value
}
