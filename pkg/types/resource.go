// pkg/types/resource.go
package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// ResourceType represents the type of infrastructure resource
type ResourceType string

// Common resource types
const (
	ResourceTypeInstance ResourceType = "instance"
	ResourceTypeVolume   ResourceType = "volume"
	ResourceTypeNetwork  ResourceType = "network"
	ResourceTypeStorage  ResourceType = "storage"
)

// ResourceStatus represents the current state of a resource
type ResourceStatus string

// Resource status constants
const (
	StatusPending     ResourceStatus = "pending"
	StatusCreating    ResourceStatus = "creating"
	StatusRunning     ResourceStatus = "running"
	StatusUpdating    ResourceStatus = "updating"
	StatusDeleting    ResourceStatus = "deleting"
	StatusDeleted     ResourceStatus = "deleted"
	StatusFailed      ResourceStatus = "failed"
	StatusUnavailable ResourceStatus = "unavailable"
)

// Resource represents any infrastructure or configuration resource
type Resource interface {
	// GetID returns the unique identifier of the resource
	GetID() string

	// GetType returns the type of the resource
	GetType() ResourceType

	// GetProvider returns the provider responsible for this resource
	GetProvider() string

	// GetStatus returns the current status of the resource
	GetStatus() ResourceStatus

	// GetMetadata returns additional resource-specific data
	GetMetadata() map[string]interface{}

	// GetTags returns the tags associated with the resource
	GetTags() map[string]string

	// GetCreatedAt returns when the resource was created
	GetCreatedAt() time.Time

	// GetUpdatedAt returns when the resource was last updated
	GetUpdatedAt() time.Time
}

// BaseResource provides a basic implementation of the Resource interface
type BaseResource struct {
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Metadata  map[string]interface{} `json:"metadata"`
	Tags      map[string]string      `json:"tags"`
	ID        string                 `json:"id"`
	Type      ResourceType           `json:"type"`
	Provider  string                 `json:"provider"`
	Status    ResourceStatus         `json:"status"`
}

// Implementation of Resource interface for BaseResource
func (r *BaseResource) GetID() string                       { return r.ID }
func (r *BaseResource) GetType() ResourceType               { return r.Type }
func (r *BaseResource) GetProvider() string                 { return r.Provider }
func (r *BaseResource) GetStatus() ResourceStatus           { return r.Status }
func (r *BaseResource) GetMetadata() map[string]interface{} { return r.Metadata }
func (r *BaseResource) GetTags() map[string]string          { return r.Tags }
func (r *BaseResource) GetCreatedAt() time.Time             { return r.CreatedAt }
func (r *BaseResource) GetUpdatedAt() time.Time             { return r.UpdatedAt }

// ResourceChange represents a change to be made to a resource
type ResourceChange struct {
	Resource     Resource
	ChangedProps map[string]interface{}
	ChangeType   ChangeType
}

// ChangeType represents the type of change to be made
type ChangeType string

const (
	ChangeTypeCreate ChangeType = "create"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
	ChangeTypeNoOp   ChangeType = "no-op"
)

// ResourceError represents an error that occurred while managing a resource
type ResourceError struct {
	Resource Resource
	Err      error
	Message  string
}

func (e *ResourceError) Error() string {
	return fmt.Sprintf("resource error [%s/%s]: %s: %v",
		e.Resource.GetProvider(),
		e.Resource.GetID(),
		e.Message,
		e.Err)
}

// ResourceDependency represents a dependency between resources
type ResourceDependency struct {
	Resource   Resource
	DependsOn  []string
	RequiredBy []string
}

// ResourceMetadata provides helper functions for working with resource metadata
type ResourceMetadata map[string]interface{}

// GetString safely retrieves a string value from metadata
func (m ResourceMetadata) GetString(key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in metadata", key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("value for key %s is not a string", key)
	}
	return s, nil
}

// GetInt safely retrieves an int value from metadata
func (m ResourceMetadata) GetInt(key string) (int, error) {
	v, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("key %s not found in metadata", key)
	}
	i, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("value for key %s is not an int", key)
	}
	return i, nil
}

// ToJSON converts the metadata to a JSON string
func (m ResourceMetadata) ToJSON() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}
	return string(bytes), nil
}

// FromJSON populates the metadata from a JSON string
func (m *ResourceMetadata) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), m)
}

// Validate checks if required metadata fields are present
func (m ResourceMetadata) Validate(required []string) error {
	for _, field := range required {
		if _, ok := m[field]; !ok {
			return fmt.Errorf("required metadata field %s is missing", field)
		}
	}
	return nil
}
