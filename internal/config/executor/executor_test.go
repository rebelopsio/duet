package executor

import (
	"context"
	"testing"
)

// mockExecutor implements the necessary methods for testing
type mockExecutor struct {
	executeFunc func(ctx context.Context, command string) (string, error)
}

func (m *mockExecutor) Execute(ctx context.Context, command string) (string, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, command)
	}
	return "", nil
}

func TestExecutor(t *testing.T) {
	t.Run("ExecuteCommand", func(t *testing.T) {
		expectedOutput := "command output"
		executor := &mockExecutor{
			executeFunc: func(ctx context.Context, command string) (string, error) {
				return expectedOutput, nil
			},
		}

		output, err := executor.Execute(context.Background(), "test command")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if output != expectedOutput {
			t.Errorf("Expected %q, got %q", expectedOutput, output)
		}
	})
}
