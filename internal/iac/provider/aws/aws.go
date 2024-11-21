package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSProvider struct {
	ec2Client *EC2Client
	region    string
}

func NewAWSProvider(ctx context.Context, region string) (*AWSProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	ec2Client, err := NewEC2Client(cfg)
	if err != nil {
		return nil, err
	}

	return &AWSProvider{
		region:    region,
		ec2Client: ec2Client,
	}, nil
}

func (p *AWSProvider) Name() string {
	return "aws"
}
