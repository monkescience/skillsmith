package registry

// EmbeddedSource is a Source backed by the embedded registry.
type EmbeddedSource struct {
	name string
}

// NewEmbeddedSource creates a new embedded source.
// The name defaults to "builtin" if empty.
func NewEmbeddedSource(name string) *EmbeddedSource {
	if name == "" {
		name = "builtin"
	}

	return &EmbeddedSource{
		name: name,
	}
}

// Name returns the source name.
func (s *EmbeddedSource) Name() string {
	return s.name
}

// Load loads all items from the embedded filesystem.
func (s *EmbeddedSource) Load() ([]Item, error) {
	reg, err := LoadFromFS(embeddedFS, "content")
	if err != nil {
		return nil, err
	}

	// Tag all items with this source
	for i := range reg.Items {
		reg.Items[i].Source = s.name
	}

	return reg.Items, nil
}
