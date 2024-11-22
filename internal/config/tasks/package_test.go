package tasks

import (
	"context"
	"testing"
)

type mockExecutor struct {
	executeFunc func(ctx context.Context, command string) (string, error)
}

func (m *mockExecutor) Execute(ctx context.Context, command string) (string, error) {
	return m.executeFunc(ctx, command)
}

func TestPackageManager(t *testing.T) {
	t.Run("InstallPackage", func(t *testing.T) {
		executorMock := &mockExecutor{
			executeFunc: func(ctx context.Context, command string) (string, error) {
				return "Package installed successfully", nil
			},
		}

		pm := &PackageManager{
			executor: executorMock,
		}

		err := pm.Install(context.Background(), "cowsay")
		if err != nil {
			t.Fatalf("Failed to install package: %v", err)
		}
	})
}
