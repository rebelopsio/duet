package ssh

import (
	"context"

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
	// Implementation
	return &Client{config: config}, nil
}

func (c *Client) Execute(ctx context.Context, command string) (string, error) {
	// Implementation
	return "", nil
}
