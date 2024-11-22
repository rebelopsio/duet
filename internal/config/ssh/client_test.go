package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"
)

// mockSSHServer simulates an SSH server for testing
type mockSSHServer struct {
	listener net.Listener
	config   *ssh.ServerConfig
}

func newMockSSHServer(t *testing.T) (*mockSSHServer, error) {
	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
	}

	privateKey, err := ssh.ParsePrivateKey([]byte(testPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	config.AddHostKey(privateKey)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &mockSSHServer{
		listener: listener,
		config:   config,
	}

	go server.serve(t)
	return server, nil
}

func (s *mockSSHServer) serve(t *testing.T) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if !isClosedError(err) {
				t.Errorf("Failed to accept connection: %v", err)
			}
			return
		}

		_, chans, reqs, err := ssh.NewServerConn(conn, s.config)
		if err != nil {
			t.Errorf("Failed to handshake: %v", err)
			continue
		}

		go ssh.DiscardRequests(reqs)
		go handleChannels(chans, t)
	}
}

func handleChannels(chans <-chan ssh.NewChannel, t *testing.T) {
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			t.Errorf("Failed to accept channel: %v", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "exec":
					payload := struct{ Command string }{}
					ssh.Unmarshal(req.Payload, &payload)

					if payload.Command == "echo test" {
						io.WriteString(channel, "test\n")
					}

					channel.Close()
				}
				req.Reply(true, nil)
			}
		}(requests)
	}
}

func isClosedError(err error) bool {
	return err.Error() == "use of closed network connection"
}

const testPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAxU4rixQXoahCL2gVoNWswNMFxYEiO0YH9YbB1qh+9nYRYGzEOc0l
...
-----END OPENSSH PRIVATE KEY-----`

func TestClient(t *testing.T) {
	mockServer, err := newMockSSHServer(t)
	if err != nil {
		t.Fatalf("Failed to start mock SSH server: %v", err)
	}
	defer mockServer.listener.Close()

	serverAddr := mockServer.listener.Addr().String()
	host, port, err := net.SplitHostPort(serverAddr)
	if err != nil {
		t.Fatalf("Failed to parse server address: %v", err)
	}

	config := &Config{
		Host:       host,
		User:       "test",
		PrivateKey: testPrivateKey,
		Port:       parseInt(port),
	}

	t.Run("Connect", func(t *testing.T) {
		client, err := NewClient(config)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		err = client.ValidateConnection()
		if err != nil {
			t.Errorf("Failed to validate connection: %v", err)
		}
	})

	t.Run("Execute", func(t *testing.T) {
		client, err := NewClient(config)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		output, err := client.Execute(context.Background(), "echo test")
		if err != nil {
			t.Fatalf("Failed to execute command: %v", err)
		}

		expected := "test\n"
		if output != expected {
			t.Errorf("Expected output %q, got %q", expected, output)
		}
	})
}

func parseInt(s string) int {
	var port int
	fmt.Sscanf(s, "%d", &port)
	return port
}
