package tasks

import (
	"context"
	"fmt"
)

// ExecutorInterface defines the required methods for command execution
type ExecutorInterface interface {
	Execute(ctx context.Context, command string) (string, error)
}

type PackageManager struct {
	executor ExecutorInterface
}

func NewPackageManager(executor ExecutorInterface) *PackageManager {
	return &PackageManager{
		executor: executor,
	}
}

func (pm *PackageManager) Install(ctx context.Context, packageName string) error {
	_, err := pm.executor.Execute(ctx, fmt.Sprintf("apt-get install -y %s", packageName))
	return err
}
