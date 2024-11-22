package planner

import (
	"context"

	"github.com/rebelopsio/duet/internal/iac/provider"
	"github.com/rebelopsio/duet/pkg/types"
)

type Change struct {
	Resource types.Resource
	Config   map[string]interface{}
	Type     string
	Provider string
}

type Plan struct {
	Changes []Change
}

type Planner struct {
	providers map[string]provider.Provider
}

func NewPlanner() *Planner {
	return &Planner{
		providers: make(map[string]provider.Provider),
	}
}

func (p *Planner) RegisterProvider(provider provider.Provider) {
	p.providers[provider.Name()] = provider
}

func (p *Planner) CreatePlan(ctx context.Context, config map[string]interface{}) (*Plan, error) {
	// Implementation will go here
	return &Plan{}, nil
}
