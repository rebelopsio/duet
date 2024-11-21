package provider

import (
	"context"
)

type Resource interface {
	ID() string
	Type() string
	Metadata() map[string]interface{}
}

type Provider interface {
	Name() string
	Create(ctx context.Context, resourceType string, config map[string]interface{}) (Resource, error)
	Read(ctx context.Context, resourceType string, id string) (Resource, error)
	Update(ctx context.Context, resource Resource, config map[string]interface{}) error
	Delete(ctx context.Context, resource Resource) error
}
