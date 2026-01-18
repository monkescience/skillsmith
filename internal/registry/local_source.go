package registry

import (
	"errors"
	"fmt"
	"os"
)

var (
	errPathNotExist = errors.New("registry path does not exist")
	errPathNotDir   = errors.New("registry path is not a directory")
)

// LocalSource is a Source backed by a local directory on disk.
type LocalSource struct {
	name string
	path string
}

// NewLocalSource creates a new local directory source.
func NewLocalSource(name, path string) *LocalSource {
	return &LocalSource{
		name: name,
		path: path,
	}
}

// Name returns the source name.
func (s *LocalSource) Name() string {
	return s.name
}

// Load loads all items from the local directory.
func (s *LocalSource) Load() ([]Item, error) {
	// Check if path exists
	info, err := os.Stat(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", errPathNotExist, s.path)
		}

		return nil, fmt.Errorf("stat registry path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%w: %s", errPathNotDir, s.path)
	}

	// Use os.DirFS to create a filesystem rooted at the path
	fsys := os.DirFS(s.path)

	reg, err := LoadFromFS(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("load from %s: %w", s.path, err)
	}

	// Tag all items with this source
	for i := range reg.Items {
		reg.Items[i].Source = s.name
	}

	return reg.Items, nil
}

// Path returns the filesystem path of this source.
func (s *LocalSource) Path() string {
	return s.path
}
