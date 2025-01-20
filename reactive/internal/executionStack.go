package internal

import "sync"

type ExecutionFrame struct {
	RegisterAsDependency func(removeDependency func(string)) (string, func())
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

func RegisterAsDependency(removeDependency func(string)) (string, func()) {
	m.Lock()
	defer m.Unlock()
	frame, err := executionStack.Peek()
	if err != nil {
		// No frame on the stack, so we don't need to add a dependency
		return "", func() {}
	}
	return frame.RegisterAsDependency(removeDependency)
}
