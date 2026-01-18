package registry

import "fmt"

// MultiRegistry aggregates items from multiple sources.
// Sources are loaded in order; first source wins for duplicate names.
type MultiRegistry struct {
	sources []Source
	items   []Item
	seen    map[string]bool // tracks item names to handle duplicates
}

// NewMultiRegistry creates a new multi-source registry.
func NewMultiRegistry() *MultiRegistry {
	return &MultiRegistry{
		sources: make([]Source, 0),
		items:   make([]Item, 0),
		seen:    make(map[string]bool),
	}
}

// AddSource adds a source to the registry.
func (m *MultiRegistry) AddSource(source Source) {
	m.sources = append(m.sources, source)
}

// Load loads items from all sources.
// First source wins for duplicate item names.
func (m *MultiRegistry) Load() error {
	m.items = make([]Item, 0)
	m.seen = make(map[string]bool)

	for _, source := range m.sources {
		items, err := source.Load()
		if err != nil {
			return fmt.Errorf("load source %s: %w", source.Name(), err)
		}

		for _, item := range items {
			// First source wins for duplicates
			if m.seen[item.Name] {
				continue
			}

			m.seen[item.Name] = true
			m.items = append(m.items, item)
		}
	}

	return nil
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
