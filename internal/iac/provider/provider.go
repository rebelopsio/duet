package provider

import (
	"context"

	"github.com/rebelopsio/duet/pkg/types"
)

// Resource is an alias for types.Resource to maintain package consistency
type Resource = types.Resource

// Provider defines the interface for infrastructure providers
type Provider interface {
	// Name returns the provider's name
	Name() string

	// Create creates a new resource
	Create(ctx context.Context, resourceType string, config map[string]interface{}) (Resource, error)

	// Read retrieves an existing resource
	Read(ctx context.Context, resourceType string, id string) (Resource, error)

	// Update updates an existing resource
	Update(ctx context.Context, resource Resource, config map[string]interface{}) error

	// Delete removes an existing resource
	Delete(ctx context.Context, resource Resource) error
}
