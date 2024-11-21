package executor

import (
	"context"
	"fmt"

	"github.com/rebelopsio/duet/internal/config/ssh"
)

type Executor struct {
	sshClient *ssh.Client
}

func NewExecutor(sshConfig *ssh.Config) (*Executor, error) {
	client, err := ssh.NewClient(sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH client: %w", err)
	}

	return &Executor{
		sshClient: client,
	}, nil
}

func (e *Executor) Execute(ctx context.Context, command string) (string, error) {
	return e.sshClient.Execute(ctx, command)
}
