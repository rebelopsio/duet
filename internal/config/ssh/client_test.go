package ssh

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

type keyPair struct {
	PublicKey  ssh.PublicKey
	PrivateKey string
}

type mockSSHServer struct {
	listener     net.Listener
	ctx          context.Context
	config       *ssh.ServerConfig
	ready        chan struct{}
	done         chan struct{}
	activeConns  map[string]net.Conn
	t            *testing.T
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	activesMutex sync.RWMutex
}

// generateTestKey generates a test RSA key pair for testing
func generateTestKey() (*keyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Convert private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Generate SSH public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	return &keyPair{
		PrivateKey: string(pem.EncodeToMemory(privateKeyPEM)),
		PublicKey:  publicKey,
	}, nil
}

func newMockSSHServer(t *testing.T, keys *keyPair) (*mockSSHServer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if key.Type() == keys.PublicKey.Type() && bytes.Equal(key.Marshal(), keys.PublicKey.Marshal()) {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("unknown public key")
		},
	}

	signer, err := ssh.ParsePrivateKey([]byte(keys.PrivateKey))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &mockSSHServer{
		listener:    listener,
		config:      config,
		ready:       make(chan struct{}),
		done:        make(chan struct{}),
		activeConns: make(map[string]net.Conn),
		t:           t,
		ctx:         ctx,
		cancel:      cancel,
	}

	server.wg.Add(1)
	go server.acceptConnections()

	// Wait for server to be ready
	select {
	case <-server.ready:
		return server, nil
	case <-time.After(2 * time.Second):
		cancel()
		if err := listener.Close(); err != nil {
			t.Logf("Failed to close listener: %v", err)
		}
		return nil, fmt.Errorf("timeout waiting for server to be ready")
	}
}

func (s *mockSSHServer) acceptConnections() {
	defer s.wg.Done()
	defer close(s.done)
	defer s.cancel()

	// Signal that we're ready to accept connections
	close(s.ready)

	for {
		if err := s.listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second)); err != nil {
			s.t.Logf("Failed to set accept deadline: %v", err)
			return
		}

		conn, err := s.listener.Accept()
		if err != nil {
			if isTimeout(err) {
				select {
				case <-s.ctx.Done():
					return
				default:
					continue
				}
			}
			if !isClosedError(err) {
				s.t.Logf("Accept error: %v", err)
			}
			return
		}

		s.activesMutex.Lock()
		s.activeConns[conn.RemoteAddr().String()] = conn
		s.activesMutex.Unlock()

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "deadline exceeded") ||
		err == context.DeadlineExceeded
}

func (s *mockSSHServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer func() {
		s.activesMutex.Lock()
		delete(s.activeConns, conn.RemoteAddr().String())
		s.activesMutex.Unlock()
		if err := conn.Close(); err != nil && !isClosedError(err) {
			s.t.Logf("Connection close error: %v", err)
		}
	}()

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, s.config)
	if err != nil {
		if !isClosedError(err) {
			s.t.Logf("SSH handshake error: %v", err)
		}
		return
	}

	go func() {
		<-s.ctx.Done()
		if err := sshConn.Close(); err != nil && !isClosedError(err) {
			s.t.Logf("SSH connection close error: %v", err)
		}
	}()

	go ssh.DiscardRequests(reqs)
	s.handleChannels(chans)
}

func (s *mockSSHServer) handleRequests(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer func() {
		if err := channel.Close(); err != nil && !isClosedError(err) {
			s.t.Logf("Channel close error: %v", err)
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case req, ok := <-requests:
			if !ok {
				return
			}
			switch req.Type {
			case "exec":
				exitStatus := make([]byte, 4)

				payload := struct{ Command string }{}
				if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
					s.t.Logf("Unmarshal error: %v", err)
					if err := req.Reply(false, nil); err != nil {
						s.t.Logf("Reply error: %v", err)
					}
					continue
				}

				if err := req.Reply(true, nil); err != nil {
					s.t.Logf("Reply error: %v", err)
					continue
				}

				if payload.Command == "echo test" {
					if _, err := io.WriteString(channel, "test\n"); err != nil {
						s.t.Logf("Write error: %v", err)
					}

					// Send exit status
					_, err := channel.SendRequest("exit-status", false, exitStatus)
					if err != nil {
						s.t.Logf("Failed to send exit status: %v", err)
					}

					if err := channel.CloseWrite(); err != nil && !isClosedError(err) {
						s.t.Logf("CloseWrite error: %v", err)
					}
					return
				}

				if payload.Command == "sleep 10" {
					// For the timeout test, we'll block until context is cancelled
					select {
					case <-s.ctx.Done():
						// Send non-zero exit status for cancelled command
						exitStatus[3] = 1
						_, err := channel.SendRequest("exit-status", false, exitStatus)
						if err != nil {
							s.t.Logf("Failed to send exit status: %v", err)
						}
					case <-time.After(10 * time.Second):
						// Normal completion
						_, err := channel.SendRequest("exit-status", false, exitStatus)
						if err != nil {
							s.t.Logf("Failed to send exit status: %v", err)
						}
					}

					if err := channel.CloseWrite(); err != nil && !isClosedError(err) {
						s.t.Logf("CloseWrite error: %v", err)
					}
					return
				}

				// Unknown command, send error exit status
				exitStatus[3] = 1
				_, err := channel.SendRequest("exit-status", false, exitStatus)
				if err != nil {
					s.t.Logf("Failed to send exit status: %v", err)
				}

				if err := channel.CloseWrite(); err != nil && !isClosedError(err) {
					s.t.Logf("CloseWrite error: %v", err)
				}
				return

			default:
				if err := req.Reply(false, nil); err != nil {
					s.t.Logf("Reply error: %v", err)
				}
			}
		}
	}
}

func (s *mockSSHServer) handleChannels(chans <-chan ssh.NewChannel) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case newChannel, ok := <-chans:
			if !ok {
				return
			}
			go s.handleChannel(newChannel)
		}
	}
}

func (s *mockSSHServer) handleChannel(newChannel ssh.NewChannel) {
	if newChannel.ChannelType() != "session" {
		if err := newChannel.Reject(ssh.UnknownChannelType, "unknown channel type"); err != nil {
			s.t.Logf("Channel reject error: %v", err)
		}
		return
	}

	channel, requests, err := newChannel.Accept()
	if err != nil {
		s.t.Logf("Channel accept error: %v", err)
		return
	}

	go s.handleRequests(channel, requests)
}

func (s *mockSSHServer) shutdown() error {
	s.cancel()

	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for server shutdown")
	}
}

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping SSH tests in short mode")
	}

	// Generate a test key for this test run
	keys, err := generateTestKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	mockServer, err := newMockSSHServer(t, keys)
	if err != nil {
		t.Fatalf("Failed to start mock SSH server: %v", err)
	}

	// Use a cleanup function to ensure proper shutdown
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- mockServer.shutdown()
		}()

		select {
		case err := <-done:
			if err != nil {
				t.Logf("Server shutdown error: %v", err)
			}
		case <-ctx.Done():
			t.Log("Server shutdown timed out")
		}
	}
	defer cleanup()

	serverAddr := mockServer.listener.Addr().String()
	host, portStr, err := net.SplitHostPort(serverAddr)
	if err != nil {
		t.Fatalf("Failed to parse server address: %v", err)
	}

	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		t.Fatalf("Failed to parse port number: %v", err)
	}

	config := &Config{
		Host:       host,
		Port:       port,
		User:       "test",
		PrivateKey: keys.PrivateKey,
		Timeout:    2 * time.Second,
	}

	// Helper function to create and cleanup client
	createClient := func(t *testing.T) (*Client, func()) {
		client, err := NewClient(config)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		cleanup := func() {
			if err := client.Close(); err != nil && !isClosedError(err) {
				t.Logf("Client close error: %v", err)
			}
		}

		return client, cleanup
	}

	t.Run("Connect", func(t *testing.T) {
		client, cleanup := createClient(t)
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- client.ValidateConnection()
		}()

		select {
		case err := <-done:
			if err != nil {
				t.Errorf("Failed to validate connection: %v", err)
			}
		case <-ctx.Done():
			t.Error("Connection validation timed out")
		}
	})

	t.Run("Execute", func(t *testing.T) {
		client, cleanup := createClient(t)
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		done := make(chan struct {
			output string
			err    error
		}, 1)

		go func() {
			output, err := client.Execute(ctx, "echo test")
			done <- struct {
				output string
				err    error
			}{output, err}
		}()

		select {
		case result := <-done:
			if result.err != nil {
				t.Fatalf("Failed to execute command: %v", result.err)
			}
			if result.output != "test\n" {
				t.Errorf("Expected output %q, got %q", "test\n", result.output)
			}
		case <-ctx.Done():
			t.Fatal("Command execution timed out")
		}
	})

	t.Run("ExecuteWithTimeout", func(t *testing.T) {
		client, cleanup := createClient(t)
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			_, err := client.Execute(ctx, "sleep 10")
			done <- err
		}()

		select {
		case err := <-done:
			if err == nil {
				t.Error("Expected error, got nil")
			} else if err != context.DeadlineExceeded {
				t.Errorf("Expected context.DeadlineExceeded error, got: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out waiting for command timeout")
		}
	})
}
