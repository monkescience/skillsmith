package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/monke/skillsmith/internal/loader"
	"github.com/monke/skillsmith/internal/registry"
	"github.com/monke/skillsmith/internal/tui"
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
	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	model := tui.NewModel(mgr)
	p := tea.NewProgram(model, tea.WithAltScreen())

	_, err = p.Run()
	if err != nil {
		return fmt.Errorf("run tui: %w", err)
	}

	return nil
}

func runList(_ *cobra.Command, _ []string) error {
	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	w := os.Stdout

	writeListOutput(w, mgr)

	return nil
}

func writeListOutput(w io.Writer, mgr *loader.Manager) {
	mustWrite(w, "Available agents and skills:\n\n")

	// List agents for all tools (use opencode as reference, items are the same)
	items := mgr.ListItems(registry.ToolOpenCode, registry.ItemTypeAgent)

	if len(items) > 0 {
		mustWrite(w, "  Agents:\n")

		for _, item := range items {
			compat := formatCompatibility(item.Compatibility)
			mustWrite(w, "    - "+item.Name+": "+item.Description+" "+compat+"\n")
		}

		mustWrite(w, "\n")
	}

	// List skills
	skills := mgr.ListItems(registry.ToolOpenCode, registry.ItemTypeSkill)

	if len(skills) > 0 {
		mustWrite(w, "  Skills:\n")

		for _, item := range skills {
			compat := formatCompatibility(item.Compatibility)
			mustWrite(w, "    - "+item.Name+": "+item.Description+" "+compat+"\n")
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
	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	registries, err := mgr.ListRegistries()
	if err != nil {
		return fmt.Errorf("list registries: %w", err)
	}

	w := os.Stdout

	mustWrite(w, "Configured registries:\n\n")

	for _, reg := range registries {
		status := ""
		if !reg.Enabled {
			status = " (disabled)"
		}

		switch reg.Type {
		case "builtin":
			mustWrite(w, fmt.Sprintf("  %s (embedded)%s\n", reg.Name, status))
		case "local":
			mustWrite(w, fmt.Sprintf("  %s: %s%s\n", reg.Name, reg.Path, status))
		case "git":
			mustWrite(w, fmt.Sprintf("  %s: %s%s\n", reg.Name, reg.URL, status))
		}
	}

	// Check if only builtin exists
	if len(registries) == 1 {
		mustWrite(w, "\nNo additional registries configured.\n")
		mustWrite(w, "Use 'skillsmith registry add <name> <path>' to add one.\n")
	}

	return nil
}

func runRegistryAdd(_ *cobra.Command, args []string) error {
	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	name := args[0]
	path := args[1]

	err = mgr.AddRegistry(name, path)
	if err != nil {
		return fmt.Errorf("add registry: %w", err)
	}

	mustWrite(os.Stdout, fmt.Sprintf("Added registry %q from %s\n", name, path))

	return nil
}

func runRegistryRemove(_ *cobra.Command, args []string) error {
	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	name := args[0]

	err = mgr.RemoveRegistry(name)
	if err != nil {
		return fmt.Errorf("remove registry: %w", err)
	}

	mustWrite(os.Stdout, fmt.Sprintf("Removed registry %q\n", name))

	return nil
}
