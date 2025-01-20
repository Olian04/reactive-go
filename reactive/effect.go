package reactive

import (
	"context"
	"time"

	"github.com/Olian04/reactive-go/reactive/internal"
	"github.com/bep/debounce"
	"github.com/google/uuid"
)

type Effect struct {
	id               string
	debounceDuration time.Duration
	fn               func()
	cancel           context.CancelFunc
}

func NewEffect(debounceDuration time.Duration, fn func()) *Effect {
	e := &Effect{
		id:               uuid.New().String(),
		debounceDuration: debounceDuration,
		fn:               fn,
	}
	e.cancel = e.run()
	return e
}

func (e *Effect) Stop() {
	if e.cancel != nil {
		e.cancel()
		e.cancel = nil
	}
}

func (e *Effect) run() context.CancelFunc {
	runDebounced := debounce.New(e.debounceDuration)
	ctx, cancel := context.WithCancel(context.Background())
	dependencyCleanupFunctions := make([]func(string), 0)

	var fun func()
	fun = func() {
		if ctx.Err() != nil {
			return
		}
		for _, fn := range dependencyCleanupFunctions {
			fn(e.id)
		}
		clear(dependencyCleanupFunctions)
		internal.PushExecutionStack(&internal.ExecutionFrame{
			RegisterAsDependency: func(removeDependency func(string)) (string, func()) {
				dependencyCleanupFunctions = append(dependencyCleanupFunctions, removeDependency)
				return e.id, func() {
					if ctx.Err() != nil {
						return
					}
					runDebounced(fun)
				}
			},
		})
		e.fn()
		internal.PopExecutionStack()
	}
	fun()
	return cancel
}
