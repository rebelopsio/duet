package ssh

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

type Config struct {
	Host       string
	User       string
	PrivateKey string
	Port       int
}

type Client struct {
	config *Config
	client *ssh.Client
}

func NewClient(config *Config) (*Client, error) {
	signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: In production, use proper host key verification
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}

	return &Client{
		config: config,
		client: client,
	}, nil
}

func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *Client) Execute(ctx context.Context, command string) (output string, err error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	defer func() {
		closeErr := session.Close()
		if err == nil {
			// If there was no error from the command, return any close error
			err = closeErr
		} else if closeErr != nil {
			// If there were both command and close errors, combine them
			err = fmt.Errorf("command error: %w; close error: %v", err, closeErr)
		}
	}()

	// Set up pipes for stdout and stderr
	var stdout, stderr io.Reader
	stdout, err = session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err = session.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := session.Start(command); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	// Read output
	outputBytes, err := io.ReadAll(stdout)
	if err != nil {
		return "", fmt.Errorf("failed to read stdout: %w", err)
	}

	// Check for errors
	errOutput, err := io.ReadAll(stderr)
	if err != nil {
		return "", fmt.Errorf("failed to read stderr: %w", err)
	}

	// Wait for the command to complete
	if err := session.Wait(); err != nil {
		return "", fmt.Errorf("command failed: %s: %w", string(errOutput), err)
	}

	return string(outputBytes), nil
}

// ValidateConnection tests the SSH connection without executing a command
func (c *Client) ValidateConnection() (err error) {
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	defer func() {
		closeErr := session.Close()
		if err == nil {
			// If validation succeeded, return any close error
			err = closeErr
		} else if closeErr != nil {
			// If both validation and close failed, combine the errors
			err = fmt.Errorf("validation error: %w; close error: %v", err, closeErr)
		}
	}()

	return nil
}
