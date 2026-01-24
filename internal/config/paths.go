package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Scope represents where to install items.
type Scope string

const (
	ScopeLocal  Scope = "local"
	ScopeGlobal Scope = "global"
)

// AllScopes returns all supported scopes.
func AllScopes() []Scope {
	return []Scope{ScopeLocal, ScopeGlobal}
}

// dirPermissions is the default permission for created directories.
const dirPermissions = 0o750

// Paths holds the resolved paths for a specific tool.
type Paths struct {
	// LocalDir is the project-local config directory.
	LocalDir string

	// GlobalDir is the user-global config directory.
	GlobalDir string

	// AgentsSubdir is the subdirectory for agents.
	AgentsSubdir string

	// SkillsSubdir is the subdirectory for skills.
	SkillsSubdir string
}

// GetPaths returns the paths for the specified tool.
// The tool parameter is a string matching registry.Tool values.
func GetPaths(tool string) (*Paths, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	switch tool {
	case "opencode":
		return &Paths{
			LocalDir:     filepath.Join(cwd, ".opencode"),
			GlobalDir:    filepath.Join(homeDir, ".config", "opencode"),
			AgentsSubdir: "agents",
			SkillsSubdir: "skills",
		}, nil

	case "claude":
		return &Paths{
			LocalDir:     filepath.Join(cwd, ".claude"),
			GlobalDir:    filepath.Join(homeDir, ".claude"),
			AgentsSubdir: "", // Claude Code doesn't have agents in the same way
			SkillsSubdir: "skills",
		}, nil

	default:
		return &Paths{
			LocalDir:     cwd,
			GlobalDir:    homeDir,
			AgentsSubdir: "agents",
			SkillsSubdir: "skills",
		}, nil
	}
}

// Exists checks if the target path already exists.
func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// EnsureDir ensures the parent directory exists.
func EnsureDir(path string) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, dirPermissions)
	if err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}

	return nil
}
