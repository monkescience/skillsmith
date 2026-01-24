package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitSource is a Source backed by a Git repository.
// The repository is cloned to a local cache directory.
type GitSource struct {
	name     string
	url      string
	cacheDir string // computed from URL hash
}

// NewGitSource creates a new Git repository source.
func NewGitSource(name, url string) *GitSource {
	return &GitSource{
		name:     name,
		url:      url,
		cacheDir: "", // computed lazily
	}
}

// Name returns the source name.
func (s *GitSource) Name() string {
	return s.name
}

// URL returns the repository URL.
func (s *GitSource) URL() string {
	return s.url
}

// Load loads all items from the Git repository.
// The repository is cloned/updated in the cache directory.
func (s *GitSource) Load() ([]Item, error) {
	cacheDir, err := s.ensureCached()
	if err != nil {
		return nil, fmt.Errorf("ensure cached: %w", err)
	}

	// Use os.DirFS to create a filesystem rooted at the cache directory
	fsys := os.DirFS(cacheDir)

	reg, err := LoadFromFS(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("load from cache: %w", err)
	}

	// Tag all items with this source
	for i := range reg.Items {
		reg.Items[i].Source = s.name
	}

	return reg.Items, nil
}

// CacheDir returns the cache directory for this source.
func (s *GitSource) CacheDir() (string, error) {
	if s.cacheDir != "" {
		return s.cacheDir, nil
	}

	baseDir, err := getCacheBaseDir()
	if err != nil {
		return "", err
	}

	// Create a hash of the URL for the directory name
	hash := sha256.Sum256([]byte(s.url))
	hashStr := hex.EncodeToString(hash[:8]) // use first 8 bytes

	// Create a safe directory name: name-hash
	safeName := sanitizeName(s.name)
	s.cacheDir = filepath.Join(baseDir, fmt.Sprintf("%s-%s", safeName, hashStr))

	return s.cacheDir, nil
}

// ensureCached ensures the repository is cloned and up-to-date.
func (s *GitSource) ensureCached() (string, error) {
	cacheDir, err := s.CacheDir()
	if err != nil {
		return "", err
	}

	// Check if cache directory exists
	if dirExists(cacheDir) {
		// Repository already cloned, try to pull latest
		err = s.pull(cacheDir)
		if err != nil {
			// Pull failed, but we can still use existing cache
			// This handles offline scenarios gracefully
			return cacheDir, nil
		}

		return cacheDir, nil
	}

	// Clone the repository
	err = s.clone(cacheDir)
	if err != nil {
		return "", fmt.Errorf("clone: %w", err)
	}

	return cacheDir, nil
}

// clone clones the repository to the specified directory.
func (s *GitSource) clone(dir string) error {
	// Ensure parent directory exists
	err := os.MkdirAll(filepath.Dir(dir), 0o755)
	if err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	// Clone with depth 1 for faster cloning
	cmd := exec.Command("git", "clone", "--depth", "1", s.url, dir)
	cmd.Stdout = nil // suppress output
	cmd.Stderr = nil

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

// pull updates the repository.
func (s *GitSource) pull(dir string) error {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = dir
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	return nil
}

// Refresh forces a fresh pull of the repository.
func (s *GitSource) Refresh() error {
	cacheDir, err := s.CacheDir()
	if err != nil {
		return err
	}

	if !dirExists(cacheDir) {
		return s.clone(cacheDir)
	}

	return s.pull(cacheDir)
}

// Clear removes the cached repository.
func (s *GitSource) Clear() error {
	cacheDir, err := s.CacheDir()
	if err != nil {
		return err
	}

	if dirExists(cacheDir) {
		return os.RemoveAll(cacheDir)
	}

	return nil
}

// getCacheBaseDir returns the base directory for caching git repositories.
func getCacheBaseDir() (string, error) {
	// Use XDG cache dir or fallback to ~/.cache
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home directory: %w", err)
		}

		cacheDir = filepath.Join(homeDir, ".cache")
	}

	return filepath.Join(cacheDir, "skillsmith", "registries"), nil
}

// sanitizeName creates a safe directory name from a source name.
func sanitizeName(name string) string {
	// Replace problematic characters
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, ":", "-")

	return name
}

// dirExists checks if a directory exists.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
