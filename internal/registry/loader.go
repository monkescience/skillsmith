package registry

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed registry/agents/*.md registry/skills/*.md
var embeddedFS embed.FS

var (
	errNoFrontmatter      = errors.New("file must start with YAML frontmatter (---)")
	errNoClosingDelimiter = errors.New("could not find closing frontmatter delimiter (---)")
)

// Load loads the embedded registry by scanning the registry directory.
func Load() (*Registry, error) {
	return LoadFromFS(embeddedFS, "registry")
}

// LoadFromFS loads a registry from a filesystem (embedded or real).
func LoadFromFS(fsys fs.FS, root string) (*Registry, error) {
	reg := &Registry{}

	// Load agents
	agentsDir := filepath.Join(root, "agents")

	agentErr := fs.WalkDir(fsys, agentsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil //nolint:nilerr // Skip on error or non-md files
		}

		item, parseErr := loadItemFromFS(fsys, path, ItemTypeAgent)
		if parseErr != nil {
			return fmt.Errorf("parse %s: %w", path, parseErr)
		}

		reg.Items = append(reg.Items, *item)

		return nil
	})
	if agentErr != nil {
		return nil, fmt.Errorf("walk agents: %w", agentErr)
	}

	// Load skills
	skillsDir := filepath.Join(root, "skills")

	skillErr := fs.WalkDir(fsys, skillsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil //nolint:nilerr // Skip on error or non-md files
		}

		item, parseErr := loadItemFromFS(fsys, path, ItemTypeSkill)
		if parseErr != nil {
			return fmt.Errorf("parse %s: %w", path, parseErr)
		}

		reg.Items = append(reg.Items, *item)

		return nil
	})
	if skillErr != nil {
		return nil, fmt.Errorf("walk skills: %w", skillErr)
	}

	return reg, nil
}

// loadItemFromFS loads a single item from a markdown file.
func loadItemFromFS(fsys fs.FS, path string, itemType ItemType) (*Item, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	item, err := ParseItem(data)
	if err != nil {
		return nil, err
	}

	item.Type = itemType
	item.SourcePath = path

	return item, nil
}

// ParseItem parses a markdown file with YAML frontmatter into an Item.
func ParseItem(data []byte) (*Item, error) {
	frontmatter, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, err
	}

	var item Item

	err = yaml.Unmarshal(frontmatter, &item)
	if err != nil {
		return nil, fmt.Errorf("parse frontmatter: %w", err)
	}

	item.Body = strings.TrimSpace(body)

	return &item, nil
}

// splitFrontmatter splits a markdown file into frontmatter and body.
func splitFrontmatter(data []byte) ([]byte, string, error) {
	const delimiter = "---"

	content := string(data)

	// Must start with ---
	if !strings.HasPrefix(content, delimiter) {
		return nil, "", errNoFrontmatter
	}

	// Find the closing ---
	rest := content[len(delimiter):]

	before, after, found := strings.Cut(rest, "\n"+delimiter)
	if !found {
		return nil, "", errNoClosingDelimiter
	}

	// Remove leading newline from body
	body := strings.TrimPrefix(after, "\n")

	return bytes.TrimSpace([]byte(before)), body, nil
}
