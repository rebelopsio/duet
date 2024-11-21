package state

import (
	"encoding/json"
	"time"
)

type ResourceState struct {
	LastUpdated   time.Time       `json:"last_updated"`
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Name          string          `json:"name"`
	Provider      string          `json:"provider"`
	Status        string          `json:"status"`
	Metadata      json.RawMessage `json:"metadata"`
	ConfigApplied bool            `json:"config_applied"`
}
