package registry

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed data/registry.yaml
var embeddedRegistry []byte

// Load loads the embedded registry.
func Load() (*Registry, error) {
	var reg Registry

	err := yaml.Unmarshal(embeddedRegistry, &reg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal registry: %w", err)
	}

	return &reg, nil
}

// LoadFromBytes loads a registry from raw YAML bytes.
func LoadFromBytes(data []byte) (*Registry, error) {
	var reg Registry

	err := yaml.Unmarshal(data, &reg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal registry: %w", err)
	}

	return &reg, nil
}
