package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Config holds the SSH client configuration
type Config struct {
	Host       string
	User       string
	PrivateKey string
	Port       int
	Timeout    time.Duration
}

// Client represents an SSH client
type Client struct {
	config *Config
	client *ssh.Client
}

// NewClient creates a new SSH client with timeouts
func NewClient(config *Config) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Parse the private key
	signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create SSH client config
	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         config.Timeout,
	}

	// Create a connection with timeout
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := net.DialTimeout("tcp", addr, config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Set connection deadline
	if err := conn.SetDeadline(time.Now().Add(config.Timeout)); err != nil {
		closeErr := conn.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to set connection deadline and close connection: %v, close error: %w", err, closeErr)
		}
		return nil, fmt.Errorf("failed to set connection deadline: %w", err)
	}

	// Create new SSH client connection
	c, chans, reqs, err := ssh.NewClientConn(conn.(*net.TCPConn), addr, sshConfig)
	if err != nil {
		closeErr := conn.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to create SSH connection and close connection: %v, close error: %w", err, closeErr)
		}
		return nil, fmt.Errorf("failed to create SSH connection: %w", err)
	}

	// Clear the deadline after successful handshake
	if err := conn.SetDeadline(time.Time{}); err != nil {
		closeErr := c.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to clear connection deadline and close client: %v, close error: %w", err, closeErr)
		}
		return nil, fmt.Errorf("failed to clear connection deadline: %w", err)
	}

	client := ssh.NewClient(c, chans, reqs)

	return &Client{
		config: config,
		client: client,
	}, nil
}

// Close closes the SSH connection
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// ValidateConnection tests if the SSH connection is working
func (c *Client) ValidateConnection() error {
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			fmt.Printf("error closing session: %v\n", err)
		}
	}()
	return nil
}

// Execute runs a command over SSH with context for cancellation
func (c *Client) Execute(ctx context.Context, command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil && !isClosedError(err) {
			fmt.Printf("error closing session: %v\n", err)
		}
	}()

	// Set up pipes for output
	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	type commandResult struct {
		output     string
		err        error
		stderrData string
	}
	resultChan := make(chan commandResult, 1)

	go func() {
		// Start the command
		if err := session.Start(command); err != nil {
			resultChan <- commandResult{err: fmt.Errorf("failed to start command: %w", err)}
			return
		}

		// Read stdout and stderr concurrently
		var stdoutData, stderrData []byte
		var stdoutErr, stderrErr error
		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			defer wg.Done()
			stdoutData, stdoutErr = io.ReadAll(stdout)
		}()

		go func() {
			defer wg.Done()
			stderrData, stderrErr = io.ReadAll(stderr)
		}()

		// Wait for all readers to complete
		wg.Wait()

		// Handle any read errors
		if stdoutErr != nil {
			resultChan <- commandResult{err: fmt.Errorf("failed to read stdout: %w", stdoutErr)}
			return
		}
		if stderrErr != nil {
			resultChan <- commandResult{err: fmt.Errorf("failed to read stderr: %w", stderrErr)}
			return
		}

		// Wait for the command to complete
		err := session.Wait()
		if err != nil {
			resultChan <- commandResult{
				err:        fmt.Errorf("command failed: %w", err),
				stderrData: string(stderrData),
			}
			return
		}

		resultChan <- commandResult{
			output:     string(stdoutData),
			stderrData: string(stderrData),
		}
	}()

	// Wait for completion or cancellation
	select {
	case result := <-resultChan:
		if result.err != nil {
			if result.stderrData != "" {
				return "", fmt.Errorf("%w: %s", result.err, result.stderrData)
			}
			return "", result.err
		}
		return result.output, nil

	case <-ctx.Done():
		if err := session.Signal(ssh.SIGTERM); err != nil && !isClosedError(err) {
			fmt.Printf("error sending SIGTERM: %v\n", err)
		}
		return "", ctx.Err()

	case <-time.After(c.config.Timeout):
		if err := session.Signal(ssh.SIGTERM); err != nil && !isClosedError(err) {
			fmt.Printf("error sending SIGTERM: %v\n", err)
		}
		return "", context.DeadlineExceeded
	}
}

func isClosedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "use of closed network connection") ||
		strings.Contains(err.Error(), "connection reset by peer") ||
		strings.Contains(err.Error(), "closed network connection") ||
		strings.Contains(err.Error(), "EOF")
}
