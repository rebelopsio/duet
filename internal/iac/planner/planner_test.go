package planner

import (
	"context"
	"testing"
	"time"

	"github.com/rebelopsio/duet/internal/iac/provider"
	"github.com/rebelopsio/duet/pkg/types"
)

// mockResource implements the types.Resource interface
type mockResource struct {
	created  time.Time
	updated  time.Time
	metadata map[string]interface{}
	tags     map[string]string
	id       string
	resType  types.ResourceType
	provider string
	status   types.ResourceStatus
}

func (m *mockResource) GetID() string                       { return m.id }
func (m *mockResource) GetType() types.ResourceType         { return m.resType }
func (m *mockResource) GetProvider() string                 { return m.provider }
func (m *mockResource) GetStatus() types.ResourceStatus     { return m.status }
func (m *mockResource) GetMetadata() map[string]interface{} { return m.metadata }
func (m *mockResource) GetTags() map[string]string          { return m.tags }
func (m *mockResource) GetCreatedAt() time.Time             { return m.created }
func (m *mockResource) GetUpdatedAt() time.Time             { return m.updated }

// mockProvider implements the provider.Provider interface
type mockProvider struct {
	name string
}

func (m *mockProvider) Name() string { return m.name }

func (m *mockProvider) Create(ctx context.Context, resourceType string, config map[string]interface{}) (provider.Resource, error) {
	return &mockResource{
		id:       "test-id",
		resType:  types.ResourceType(resourceType),
		provider: m.name,
		status:   types.StatusRunning,
		metadata: config,
		tags:     make(map[string]string),
		created:  time.Now(),
		updated:  time.Now(),
	}, nil
}

func (m *mockProvider) Read(ctx context.Context, resourceType string, id string) (provider.Resource, error) {
	return nil, nil
}

func (m *mockProvider) Update(ctx context.Context, resource provider.Resource, config map[string]interface{}) error {
	return nil
}

func (m *mockProvider) Delete(ctx context.Context, resource provider.Resource) error {
	return nil
}

func TestPlanner(t *testing.T) {
	t.Run("CreatePlan", func(t *testing.T) {
		planner := &Planner{
			providers: make(map[string]provider.Provider),
		}

		mockAWSProvider := &mockProvider{name: "aws"}
		planner.RegisterProvider(mockAWSProvider)

		config := map[string]interface{}{
			"provider": "aws",
			"resources": []map[string]interface{}{
				{
					"type": "instance",
					"name": "test-instance",
				},
			},
		}

		plan, err := planner.CreatePlan(context.Background(), config)
		if err != nil {
			t.Fatalf("Failed to create plan: %v", err)
		}

		if plan == nil {
			t.Error("Expected plan to not be nil")
		}
	})
}
