package registry

import (
	"fmt"
	"strings"
)

// MultiRegistry aggregates items from multiple sources.
// Sources are loaded in order; last source wins for duplicate names.
// This allows project-specific registries to override builtin/global items.
type MultiRegistry struct {
	sources  []Source
	items    []Item
	itemsMap map[string]int // maps item name to index in items slice
	errors   []error        // errors from sources that failed to load
}

// NewMultiRegistry creates a new multi-source registry.
func NewMultiRegistry() *MultiRegistry {
	return &MultiRegistry{
		sources:  make([]Source, 0),
		items:    make([]Item, 0),
		itemsMap: make(map[string]int),
		errors:   make([]error, 0),
	}
}

// AddSource adds a source to the registry.
func (m *MultiRegistry) AddSource(source Source) {
	m.sources = append(m.sources, source)
}

// Load loads items from all sources.
// Last source wins for duplicate item names, allowing overrides.
// If a source fails to load, its error is recorded but other sources continue loading.
func (m *MultiRegistry) Load() error {
	m.items = make([]Item, 0)
	m.itemsMap = make(map[string]int)
	m.errors = make([]error, 0)

	for _, source := range m.sources {
		items, err := source.Load()
		if err != nil {
			// Record the error but continue with other sources
			m.errors = append(m.errors, fmt.Errorf("source %s: %w", source.Name(), err))

			continue
		}

		for _, item := range items {
			// Last source wins - replace existing item if present
			if idx, exists := m.itemsMap[item.Name]; exists {
				m.items[idx] = item
			} else {
				m.itemsMap[item.Name] = len(m.items)
				m.items = append(m.items, item)
			}
		}
	}

	return nil
}

// Errors returns any errors that occurred during loading.
func (m *MultiRegistry) Errors() []error {
	return m.errors
}

// HasErrors returns true if any sources failed to load.
func (m *MultiRegistry) HasErrors() bool {
	return len(m.errors) > 0
}

// ErrorString returns a combined error message for all loading errors.
func (m *MultiRegistry) ErrorString() string {
	if len(m.errors) == 0 {
		return ""
	}

	msgs := make([]string, len(m.errors))
	for i, err := range m.errors {
		msgs[i] = err.Error()
	}

	return strings.Join(msgs, "; ")
}

// Registry returns the aggregated registry.
func (m *MultiRegistry) Registry() *Registry {
	return &Registry{Items: m.items}
}

// Sources returns the list of configured sources.
func (m *MultiRegistry) Sources() []Source {
	return m.sources
}

// LoadWithBuiltin creates a MultiRegistry with only the builtin embedded source.
func LoadWithBuiltin() (*MultiRegistry, error) {
	multi := NewMultiRegistry()
	multi.AddSource(NewEmbeddedSource("builtin"))

	err := multi.Load()
	if err != nil {
		return nil, err
	}

	return multi, nil
}
