package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/loader"
	"github.com/monke/skillsmith/internal/registry"
	"github.com/monke/skillsmith/internal/tui"
)

var (
	errPathNotExist        = errors.New("path does not exist")
	errPathNotDir          = errors.New("path is not a directory")
	errRegistryExists      = errors.New("registry with this name already exists")
	errCannotRemoveBuiltin = errors.New("cannot remove the builtin registry")
	errRegistryNotFound    = errors.New("registry not found")
)

var version = "dev"

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "skillsmith",
	Short: "Install agents and skills for AI coding tools",
	Long: `skillsmith is a TUI for browsing, previewing, and installing
agents, subagents, and skills for AI coding tools like OpenCode and Claude Code.

Run 'skillsmith tui' to launch the interactive browser.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI browser",
	Long: `Launch the interactive TUI to browse and install agents and skills.

Use arrow keys or vim bindings (h/j/k/l) to navigate.
Press Enter to install locally, 'g' to install globally.
Press '?' for help.`,
	RunE: runTUI,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available agents and skills",
	RunE:  runList,
}

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage registry sources",
	Long:  `Manage the registry sources that skillsmith loads agents and skills from.`,
}

var registryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured registries",
	RunE:  runRegistryList,
}

var registryAddCmd = &cobra.Command{
	Use:   "add <name> <path>",
	Short: "Add a local registry",
	Long: `Add a local directory as a registry source.

The directory should contain 'agents/' and/or 'skills/' subdirectories
with markdown files using YAML frontmatter.`,
	Args: cobra.ExactArgs(2), //nolint:mnd // name and path
	RunE: runRegistryAdd,
}

var registryRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a registry",
	Args:  cobra.ExactArgs(1),
	RunE:  runRegistryRemove,
}

func setupCommands() {
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(registryCmd)

	registryCmd.AddCommand(registryListCmd)
	registryCmd.AddCommand(registryAddCmd)
	registryCmd.AddCommand(registryRemoveCmd)
}

//nolint:gochecknoinits // cobra requires init for command setup
func init() {
	setupCommands()
}

func runTUI(_ *cobra.Command, _ []string) error {
	multi, err := loader.LoadFromConfig()
	if err != nil {
		return fmt.Errorf("load registry: %w", err)
	}

	model := tui.NewModel(multi.Registry())
	p := tea.NewProgram(model, tea.WithAltScreen())

	_, err = p.Run()
	if err != nil {
		return fmt.Errorf("run tui: %w", err)
	}

	return nil
}

func runList(_ *cobra.Command, _ []string) error {
	multi, err := loader.LoadFromConfig()
	if err != nil {
		return fmt.Errorf("load registry: %w", err)
	}

	w := os.Stdout

	writeListOutput(w, multi.Registry())

	return nil
}

func writeListOutput(w io.Writer, reg *registry.Registry) {
	mustWrite(w, "Available agents and skills:\n\n")

	// List agents
	agents := reg.ByType(registry.ItemTypeAgent)
	if len(agents) > 0 {
		mustWrite(w, "  Agents:\n")

		for _, agent := range agents {
			compat := formatCompatibility(agent.Compatibility)
			mustWrite(w, "    - "+agent.Name+": "+agent.Description+" "+compat+"\n")
		}

		mustWrite(w, "\n")
	}

	// List skills
	skills := reg.ByType(registry.ItemTypeSkill)
	if len(skills) > 0 {
		mustWrite(w, "  Skills:\n")

		for _, skill := range skills {
			compat := formatCompatibility(skill.Compatibility)
			mustWrite(w, "    - "+skill.Name+": "+skill.Description+" "+compat+"\n")
		}

		mustWrite(w, "\n")
	}
}

func formatCompatibility(tools []registry.Tool) string {
	if len(tools) == 0 {
		return ""
	}

	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = string(t)
	}

	return "[" + strings.Join(names, ", ") + "]"
}

func mustWrite(w io.Writer, s string) {
	_, _ = w.Write([]byte(s))
}

func runRegistryList(_ *cobra.Command, _ []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	w := os.Stdout

	mustWrite(w, "Configured registries:\n\n")
	mustWrite(w, "  builtin (embedded)\n")

	if len(cfg.Registries) == 0 {
		mustWrite(w, "\nNo additional registries configured.\n")
		mustWrite(w, "Use 'skillsmith registry add <name> <path>' to add one.\n")

		return nil
	}

	for _, reg := range cfg.Registries {
		status := ""
		if !reg.IsEnabled() {
			status = " (disabled)"
		}

		switch {
		case reg.IsLocal():
			mustWrite(w, fmt.Sprintf("  %s: %s%s\n", reg.Name, reg.Path, status))
		case reg.IsGit():
			mustWrite(w, fmt.Sprintf("  %s: %s%s\n", reg.Name, reg.URL, status))
		}
	}

	return nil
}

func runRegistryAdd(_ *cobra.Command, args []string) error {
	name := args[0]
	path := args[1]

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home directory: %w", err)
		}

		path = filepath.Join(homeDir, path[2:])
	}

	// Make path absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	// Verify path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", errPathNotExist, absPath)
		}

		return fmt.Errorf("check path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: %s", errPathNotDir, absPath)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Check for duplicate name
	for _, reg := range cfg.Registries {
		if reg.Name == name {
			return fmt.Errorf("%w: %s", errRegistryExists, name)
		}
	}

	cfg.Registries = append(cfg.Registries, config.RegistrySource{
		Name: name,
		Path: absPath,
	})

	err = config.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	mustWrite(os.Stdout, fmt.Sprintf("Added registry %q from %s\n", name, absPath))

	return nil
}

func runRegistryRemove(_ *cobra.Command, args []string) error {
	name := args[0]

	if name == "builtin" {
		return errCannotRemoveBuiltin
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	found := false
	newRegistries := make([]config.RegistrySource, 0, len(cfg.Registries))

	for _, reg := range cfg.Registries {
		if reg.Name == name {
			found = true

			continue
		}

		newRegistries = append(newRegistries, reg)
	}

	if !found {
		return fmt.Errorf("%w: %s", errRegistryNotFound, name)
	}

	cfg.Registries = newRegistries

	err = config.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	mustWrite(os.Stdout, fmt.Sprintf("Removed registry %q\n", name))

	return nil
}
