package tasks

import (
	"context"

	"github.com/rebelopsio/duet/internal/config/executor"
)

type PackageManager struct {
	executor *executor.Executor
}

func NewPackageManager(executor *executor.Executor) *PackageManager {
	return &PackageManager{
		executor: executor,
	}
}

func (pm *PackageManager) Install(ctx context.Context, packageName string) error {
	// Implementation
	return nil
}
