package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monke/skillsmith/internal/registry"
)

// Scope represents where to install items.
type Scope string

const (
	ScopeLocal  Scope = "local"
	ScopeGlobal Scope = "global"
)

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
func GetPaths(tool registry.Tool) (*Paths, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	switch tool {
	case registry.ToolOpenCode:
		return &Paths{
			LocalDir:     filepath.Join(cwd, ".opencode"),
			GlobalDir:    filepath.Join(homeDir, ".config", "opencode"),
			AgentsSubdir: "agents",
			SkillsSubdir: "skills",
		}, nil

	case registry.ToolClaudeCode:
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

// GetInstallPath returns the full path where an item should be installed.
func GetInstallPath(item registry.Item, scope Scope) (string, error) {
	paths, err := GetPaths(item.Tool)
	if err != nil {
		return "", err
	}

	var baseDir string
	if scope == ScopeGlobal {
		baseDir = paths.GlobalDir
	} else {
		baseDir = paths.LocalDir
	}

	switch item.Type {
	case registry.ItemTypeAgent:
		if paths.AgentsSubdir == "" {
			// For tools without agent subdirs, install directly
			return filepath.Join(baseDir, item.Filename), nil
		}

		return filepath.Join(baseDir, paths.AgentsSubdir, item.Filename), nil

	case registry.ItemTypeSkill:
		// Skills go in skills/<name>/SKILL.md
		skillDir := filepath.Join(baseDir, paths.SkillsSubdir, item.Name)

		return filepath.Join(skillDir, "SKILL.md"), nil

	default:
		return filepath.Join(baseDir, item.Filename), nil
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
