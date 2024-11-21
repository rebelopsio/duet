// internal/iac/provider/aws/ec2.go
package aws

import (
	"context"
	"fmt"

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
	// Example implementation using the client
	input := &ec2.RunInstancesInput{
		MaxCount: aws.Int32(1),
		MinCount: aws.Int32(1),
		// Add other configuration as needed
	}

	result, err := c.client.RunInstances(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create instance: %w", err)
	}

	if len(result.Instances) == 0 {
		return "", fmt.Errorf("no instance created")
	}

	return *result.Instances[0].InstanceId, nil
}
