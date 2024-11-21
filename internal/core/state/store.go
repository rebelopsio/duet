// internal/core/state/store.go
package state

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

type Resource struct {
	ID            string `gorm:"primaryKey"`
	Type          string
	Name          string
	Provider      string
	Status        string
	LastUpdated   string
	Metadata      []byte
	ConfigApplied bool
}

func NewStore(dbPath string) (*Store, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.AutoMigrate(&Resource{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Store{db: db}, nil
}

// GetResources retrieves all resources from the store
func (s *Store) GetResources(ctx context.Context) ([]Resource, error) {
	var resources []Resource
	result := s.db.WithContext(ctx).Find(&resources)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get resources: %w", result.Error)
	}
	return resources, nil
}

// SaveResource saves a resource to the store
func (s *Store) SaveResource(ctx context.Context, resource *Resource) error {
	result := s.db.WithContext(ctx).Save(resource)
	return result.Error
}

// GetResource retrieves a single resource by ID
func (s *Store) GetResource(ctx context.Context, id string) (*Resource, error) {
	var resource Resource
	result := s.db.WithContext(ctx).First(&resource, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resource, nil
}

// DeleteResource removes a resource from the store
func (s *Store) DeleteResource(ctx context.Context, id string) error {
	result := s.db.WithContext(ctx).Delete(&Resource{}, "id = ?", id)
	return result.Error
}
