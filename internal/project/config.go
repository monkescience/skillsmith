// Package project provides project-specific configuration handling.
package project

import (
	"slices"

	"github.com/monke/skillsmith/internal/config"
)

// ConfigFileName is the name of the project configuration file.
const ConfigFileName = ".skillsmith.yaml"

// Config represents a project's skillsmith configuration.
// This file lives in the project root and defines which skills/agents
// should be installed for the project.
type Config struct {
	// Tools limits installation to specific tools.
	// If empty, skills are installed for ALL compatible tools.
	// Valid values: "claude", "opencode"
	Tools []string `yaml:"tools,omitempty"`

	// Registries are project-specific registry sources.
	// These are merged with global registries, with project registries
	// taking precedence (last wins for duplicate skill names).
	Registries []config.RegistrySource `yaml:"registries,omitempty"`

	// Skills lists the skills to install for this project.
	Skills []string `yaml:"skills,omitempty"`

	// Agents lists the agents to install for this project.
	Agents []string `yaml:"agents,omitempty"`
}

// HasTool returns true if the config includes the specified tool,
// or if no tools are specified (meaning all tools).
func (c *Config) HasTool(tool string) bool {
	if len(c.Tools) == 0 {
		return true
	}

	return slices.Contains(c.Tools, tool)
}

// AllItems returns all skills and agents combined.
func (c *Config) AllItems() []string {
	items := make([]string, 0, len(c.Skills)+len(c.Agents))
	items = append(items, c.Skills...)
	items = append(items, c.Agents...)

	return items
}

// IsEmpty returns true if the config has no skills or agents defined.
func (c *Config) IsEmpty() bool {
	return len(c.Skills) == 0 && len(c.Agents) == 0
}

// AddSkill adds a skill to the config if not already present.
// Returns true if the skill was added, false if it already existed.
func (c *Config) AddSkill(name string) bool {
	if slices.Contains(c.Skills, name) {
		return false
	}

	c.Skills = append(c.Skills, name)

	return true
}

// RemoveSkill removes a skill from the config.
// Returns true if the skill was removed, false if it wasn't present.
func (c *Config) RemoveSkill(name string) bool {
	for i, s := range c.Skills {
		if s == name {
			c.Skills = append(c.Skills[:i], c.Skills[i+1:]...)
			return true
		}
	}

	return false
}

// AddAgent adds an agent to the config if not already present.
// Returns true if the agent was added, false if it already existed.
func (c *Config) AddAgent(name string) bool {
	if slices.Contains(c.Agents, name) {
		return false
	}

	c.Agents = append(c.Agents, name)

	return true
}

// RemoveAgent removes an agent from the config.
// Returns true if the agent was removed, false if it wasn't present.
func (c *Config) RemoveAgent(name string) bool {
	for i, a := range c.Agents {
		if a == name {
			c.Agents = append(c.Agents[:i], c.Agents[i+1:]...)
			return true
		}
	}

	return false
}

// HasSkill returns true if the skill is in the config.
func (c *Config) HasSkill(name string) bool {
	for _, s := range c.Skills {
		if s == name {
			return true
		}
	}

	return false
}

// HasAgent returns true if the agent is in the config.
func (c *Config) HasAgent(name string) bool {
	for _, a := range c.Agents {
		if a == name {
			return true
		}
	}

	return false
}
