package state

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestStore(t *testing.T) {
	// Create temporary database file
	tmpDir, err := os.MkdirTemp("", "duet-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// Initialize store
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	ctx := context.Background()

	t.Run("SaveAndGetResource", func(t *testing.T) {
		resource := &Resource{
			ID:            "test-resource-1",
			Type:          "test",
			Name:          "Test Resource",
			Provider:      "test-provider",
			Status:        "running",
			LastUpdated:   "2024-01-01",
			ConfigApplied: true,
		}

		// Save resource
		err := store.SaveResource(ctx, resource)
		if err != nil {
			t.Fatalf("Failed to save resource: %v", err)
		}

		// Get resource
		retrieved, err := store.GetResource(ctx, resource.ID)
		if err != nil {
			t.Fatalf("Failed to get resource: %v", err)
		}

		if retrieved.ID != resource.ID {
			t.Errorf("Expected resource ID %s, got %s", resource.ID, retrieved.ID)
		}
		if retrieved.Type != resource.Type {
			t.Errorf("Expected resource type %s, got %s", resource.Type, retrieved.Type)
		}
	})

	t.Run("GetResources", func(t *testing.T) {
		// Add another resource
		resource2 := &Resource{
			ID:            "test-resource-2",
			Type:          "test",
			Name:          "Test Resource 2",
			Provider:      "test-provider",
			Status:        "running",
			LastUpdated:   "2024-01-01",
			ConfigApplied: true,
		}

		err := store.SaveResource(ctx, resource2)
		if err != nil {
			t.Fatalf("Failed to save second resource: %v", err)
		}

		// Get all resources
		resources, err := store.GetResources(ctx)
		if err != nil {
			t.Fatalf("Failed to get resources: %v", err)
		}

		if len(resources) != 2 {
			t.Errorf("Expected 2 resources, got %d", len(resources))
		}
	})

	t.Run("DeleteResource", func(t *testing.T) {
		// Delete a resource
		err := store.DeleteResource(ctx, "test-resource-1")
		if err != nil {
			t.Fatalf("Failed to delete resource: %v", err)
		}

		// Verify deletion
		_, err = store.GetResource(ctx, "test-resource-1")
		if err == nil {
			t.Error("Expected error getting deleted resource, got nil")
		}

		// Check remaining resources
		resources, err := store.GetResources(ctx)
		if err != nil {
			t.Fatalf("Failed to get resources: %v", err)
		}

		if len(resources) != 1 {
			t.Errorf("Expected 1 resource after deletion, got %d", len(resources))
		}
	})
}
