package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Client struct {
	client *ec2.Client
}

func NewEC2Client(cfg aws.Config) (*EC2Client, error) {
	return &EC2Client{
		client: ec2.NewFromConfig(cfg),
	}, nil
}

func (c *EC2Client) CreateInstance(ctx context.Context, config map[string]interface{}) (string, error) {
	// Implementation
	return "", nil
}
