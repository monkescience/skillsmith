package registry

import "io/fs"

// Source represents a registry source that provides items.
// Sources can be embedded, local directories, or remote Git repositories.
type Source interface {
	// Name returns the unique name of this source.
	Name() string

	// Load loads all items from this source.
	// Items should have their Source field set to identify origin.
	Load() ([]Item, error)
}

// FSSource is a Source backed by a filesystem (embedded or real).
type FSSource struct {
	name string
	fs   fs.FS
	root string
}

// NewFSSource creates a new filesystem-backed source.
func NewFSSource(name string, fsys fs.FS, root string) *FSSource {
	return &FSSource{
		name: name,
		fs:   fsys,
		root: root,
	}
}

// Name returns the source name.
func (s *FSSource) Name() string {
	return s.name
}

// Load loads all items from the filesystem.
func (s *FSSource) Load() ([]Item, error) {
	reg, err := LoadFromFS(s.fs, s.root)
	if err != nil {
		return nil, err
	}

	// Tag all items with this source
	for i := range reg.Items {
		reg.Items[i].Source = s.name
	}

	return reg.Items, nil
}
