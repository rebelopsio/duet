package planner

import (
	"context"

	"github.com/rebelopsio/duet/internal/core/state"
	"github.com/rebelopsio/duet/internal/iac/provider"
)

type Change struct {
	Resource provider.Resource
	Config   map[string]interface{}
	Type     string
	Provider string
}

type Plan struct {
	Changes []Change
}

type Planner struct {
	store     *state.Store
	providers map[string]provider.Provider
}

func NewPlanner(store *state.Store) *Planner {
	return &Planner{
		store:     store,
		providers: make(map[string]provider.Provider),
	}
}

func (p *Planner) RegisterProvider(provider provider.Provider) {
	p.providers[provider.Name()] = provider
}

func (p *Planner) CreatePlan(ctx context.Context, config map[string]interface{}) (*Plan, error) {
	// Implementation
	return &Plan{}, nil
}
