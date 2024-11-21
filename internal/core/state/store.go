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

func NewStore(dbPath string) (*Store, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.AutoMigrate(&ResourceState{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) SaveResource(ctx context.Context, resource *ResourceState) error {
	result := s.db.WithContext(ctx).Save(resource)
	return result.Error
}

func (s *Store) GetResource(ctx context.Context, id string) (*ResourceState, error) {
	var resource ResourceState
	result := s.db.WithContext(ctx).First(&resource, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resource, nil
}
