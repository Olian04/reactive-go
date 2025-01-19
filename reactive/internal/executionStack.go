package internal

import "sync"

type ExecutionFrame struct {
	AddDependency func(id string, fn func() bool) (string, func())
}

var m sync.Mutex
var executionStack = NewStack[*ExecutionFrame]()

func PushExecutionStack(frame *ExecutionFrame) {
	m.Lock()
	defer m.Unlock()
	executionStack.Push(frame)
}

func PopExecutionStack() *ExecutionFrame {
	m.Lock()
	defer m.Unlock()
	frame, err := executionStack.Pop()
	if err != nil {
		panic(err)
	}
	return frame
}

func AddDependency(id string, fn func() bool) (string, func()) {
	m.Lock()
	defer m.Unlock()
	frame, err := executionStack.Peek()
	if err != nil {
		// No frame on the stack, so we don't need to add a dependency
		return "", func() {}
	}
	return frame.AddDependency(id, fn)
}
