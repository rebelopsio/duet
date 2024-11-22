package executor

import (
	"context"
)

// Executor defines the interface for executing commands
type ExecutorInterface interface {
	Execute(ctx context.Context, command string) (string, error)
}

type Executor struct {
	executor ExecutorInterface
}

func NewExecutor(executor ExecutorInterface) *Executor {
	return &Executor{
		executor: executor,
	}
}

func (e *Executor) Execute(ctx context.Context, command string) (string, error) {
	return e.executor.Execute(ctx, command)
}
