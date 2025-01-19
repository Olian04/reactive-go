package internal

type ExecutionFrame struct {
	AddDependency func(id string, fn func() bool) (string, func())
}

var executionStack = NewStack[*ExecutionFrame]()

func PushExecutionStack(frame *ExecutionFrame) {
	executionStack.Push(frame)
}

func PopExecutionStack() *ExecutionFrame {
	frame, err := executionStack.Pop()
	if err != nil {
		panic(err)
	}
	return frame
}

func AddDependency(id string, fn func() bool) (string, func()) {
	frame, err := executionStack.Peek()
	if err != nil {
		// No frame on the stack, so we don't need to add a dependency
		return "", func() {}
	}
	return frame.AddDependency(id, fn)
}
