package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/monke/skillsmith/internal/config"
	"github.com/monke/skillsmith/internal/loader"
	"github.com/monke/skillsmith/internal/project"
	"github.com/monke/skillsmith/internal/registry"
	"github.com/monke/skillsmith/internal/tui"
)

// Command errors.
var (
	errProjectExists    = errors.New("project already exists")
	errNoProject        = errors.New("no project found, run 'skillsmith project init' first")
	errItemNotFound     = errors.New("item not found in registry")
	errUnknownItemType  = errors.New("unknown item type")
	errItemNotInProject = errors.New("item is not in the project")
	errInstallFailed    = errors.New("some items failed to install")
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

var registryAddGitCmd = &cobra.Command{
	Use:   "add-git <name> <url>",
	Short: "Add a Git registry",
	Long: `Add a Git repository as a registry source.

The repository will be cloned to a local cache and used as a source
for agents and skills. The repository should contain 'agents/' and/or
'skills/' subdirectories with markdown files using YAML frontmatter.

Example:
  skillsmith registry add-git team-skills https://github.com/myteam/skills.git`,
	Args: cobra.ExactArgs(2), //nolint:mnd // name and url
	RunE: runRegistryAddGit,
}

// Project commands.
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project-specific skills",
	Long: `Manage project-specific skills and agents.

A project is defined by a .skillsmith.yaml file in the project root.
This file lists the skills and agents to install for the project.`,
}

var projectInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Long: `Create a new .skillsmith.yaml file in the current directory.

This file defines which skills and agents should be installed for this project.`,
	RunE: runProjectInit,
}

var projectAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a skill or agent to the project",
	Long: `Add a skill or agent to the project's .skillsmith.yaml file.

The item will be added to either the skills or agents list based on its type.`,
	Args: cobra.ExactArgs(1),
	RunE: runProjectAdd,
}

var projectRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a skill or agent from the project",
	Args:  cobra.ExactArgs(1),
	RunE:  runProjectRemove,
}

var projectInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install all project skills and agents",
	Long: `Install all skills and agents defined in .skillsmith.yaml.

Items are installed for all compatible tools, or only for tools specified
in the project config.`,
	RunE: runProjectInstall,
}

var projectStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show project installation status",
	Long:  `Show the installation status of all skills and agents in the project.`,
	RunE:  runProjectStatus,
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List skills and agents in the project",
	RunE:  runProjectList,
}

// Flags.
var projectInstallForce bool

func setupCommands() {
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(registryCmd)
	rootCmd.AddCommand(projectCmd)

	registryCmd.AddCommand(registryListCmd)
	registryCmd.AddCommand(registryAddCmd)
	registryCmd.AddCommand(registryRemoveCmd)
	registryCmd.AddCommand(registryAddGitCmd)

	projectCmd.AddCommand(projectInitCmd)
	projectCmd.AddCommand(projectAddCmd)
	projectCmd.AddCommand(projectRemoveCmd)
	projectCmd.AddCommand(projectInstallCmd)
	projectCmd.AddCommand(projectStatusCmd)
	projectCmd.AddCommand(projectListCmd)

	// Flags
	projectInstallCmd.Flags().BoolVarP(&projectInstallForce, "force", "f", false, "Force reinstall even if up to date")
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

func runRegistryAddGit(_ *cobra.Command, args []string) error {
	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	name := args[0]
	url := args[1]

	err = mgr.AddGitRegistry(name, url)
	if err != nil {
		return fmt.Errorf("add git registry: %w", err)
	}

	mustWrite(os.Stdout, fmt.Sprintf("Added git registry %q from %s\n", name, url))
	mustWrite(os.Stdout, "The repository will be cloned when you next load skills.\n")

	return nil
}

// Project command implementations.

func runProjectInit(_ *cobra.Command, _ []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Check if project already exists
	if project.ExistsInDir(cwd) {
		return fmt.Errorf("%w: %s", errProjectExists, project.GetConfigPath(cwd))
	}

	cfg, err := project.Init(cwd)
	if err != nil {
		return fmt.Errorf("initialize project: %w", err)
	}

	_ = cfg // unused for now

	mustWrite(os.Stdout, fmt.Sprintf("Created %s\n", project.GetConfigPath(cwd)))
	mustWrite(os.Stdout, "\nNext steps:\n")
	mustWrite(os.Stdout, "  skillsmith project add <skill>   Add skills to the project\n")
	mustWrite(os.Stdout, "  skillsmith project install       Install all project skills\n")

	return nil
}

func runProjectAdd(_ *cobra.Command, args []string) error {
	name := args[0]

	// Load project config
	cfg, projectDir, err := project.Load()
	if err != nil {
		if errors.Is(err, project.ErrNotFound) {
			return errNoProject
		}

		return fmt.Errorf("load project: %w", err)
	}

	// Load manager to check if item exists
	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	// Find the item to determine its type
	item, err := mgr.GetItem(name)
	if err != nil {
		return fmt.Errorf("%w: %s", errItemNotFound, name)
	}

	// Add to appropriate list
	var added bool

	switch item.Type {
	case registry.ItemTypeSkill:
		added = cfg.AddSkill(name)
	case registry.ItemTypeAgent:
		added = cfg.AddAgent(name)
	default:
		return fmt.Errorf("%w: %s", errUnknownItemType, item.Type)
	}

	if !added {
		mustWrite(os.Stdout, fmt.Sprintf("%s %q is already in the project\n", item.Type, name))

		return nil
	}

	// Save config
	err = project.Save(cfg, projectDir)
	if err != nil {
		return fmt.Errorf("save project: %w", err)
	}

	mustWrite(os.Stdout, fmt.Sprintf("Added %s %q to project\n", item.Type, name))

	return nil
}

func runProjectRemove(_ *cobra.Command, args []string) error {
	name := args[0]

	cfg, projectDir, err := project.Load()
	if err != nil {
		if errors.Is(err, project.ErrNotFound) {
			return errNoProject
		}

		return fmt.Errorf("load project: %w", err)
	}

	// Try to remove from both lists
	removedSkill := cfg.RemoveSkill(name)
	removedAgent := cfg.RemoveAgent(name)

	if !removedSkill && !removedAgent {
		return fmt.Errorf("%w: %s", errItemNotInProject, name)
	}

	err = project.Save(cfg, projectDir)
	if err != nil {
		return fmt.Errorf("save project: %w", err)
	}

	mustWrite(os.Stdout, fmt.Sprintf("Removed %q from project\n", name))

	return nil
}

func runProjectInstall(_ *cobra.Command, _ []string) error {
	// Load project config
	cfg, _, err := project.Load()
	if err != nil {
		if errors.Is(err, project.ErrNotFound) {
			return errNoProject
		}

		return fmt.Errorf("load project: %w", err)
	}

	if cfg.IsEmpty() {
		mustWrite(os.Stdout, "No skills or agents defined in project.\n")
		mustWrite(os.Stdout, "Use 'skillsmith project add <name>' to add items.\n")

		return nil
	}

	// Load manager
	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	// Install all items
	results := mgr.InstallProjectItems(cfg, config.ScopeLocal, projectInstallForce)

	// Display results
	w := os.Stdout

	var installed, skipped, failed int

	for _, r := range results {
		if r.Error != nil {
			failed++
			mustWrite(w, fmt.Sprintf("  [FAIL] %s (%s): %v\n", r.ItemName, r.Tool, r.Error))
		} else if r.Skipped {
			skipped++

			if r.Success {
				// Already installed/up to date
				mustWrite(w, fmt.Sprintf("  [OK]   %s (%s): %s\n", r.ItemName, r.Tool, r.Reason))
			} else {
				mustWrite(w, fmt.Sprintf("  [SKIP] %s (%s): %s\n", r.ItemName, r.Tool, r.Reason))
			}
		} else if r.Success {
			installed++
			mustWrite(w, fmt.Sprintf("  [NEW]  %s (%s) -> %s\n", r.ItemName, r.Tool, r.Path))
		}
	}

	mustWrite(w, fmt.Sprintf("\nInstalled: %d, Skipped: %d, Failed: %d\n", installed, skipped, failed))

	if failed > 0 {
		return errInstallFailed
	}

	return nil
}

func runProjectStatus(_ *cobra.Command, _ []string) error {
	cfg, _, err := project.Load()
	if err != nil {
		if errors.Is(err, project.ErrNotFound) {
			return errNoProject
		}

		return fmt.Errorf("load project: %w", err)
	}

	if cfg.IsEmpty() {
		mustWrite(os.Stdout, "No skills or agents defined in project.\n")

		return nil
	}

	mgr, err := loader.NewManager()
	if err != nil {
		return fmt.Errorf("initialize manager: %w", err)
	}

	results := mgr.GetProjectStatus(cfg, config.ScopeLocal)

	w := os.Stdout

	mustWrite(w, "Project status:\n\n")

	for _, r := range results {
		status := "[ ]"
		if r.Success {
			status = "[x]"
		}

		if r.Error != nil {
			mustWrite(w, fmt.Sprintf("  [!] %s (%s): %v\n", r.ItemName, r.Tool, r.Error))
		} else if r.Skipped {
			mustWrite(w, fmt.Sprintf("  [-] %s (%s): %s\n", r.ItemName, r.Tool, r.Reason))
		} else {
			mustWrite(w, fmt.Sprintf("  %s %s (%s): %s\n", status, r.ItemName, r.Tool, r.Reason))
		}
	}

	return nil
}

func runProjectList(_ *cobra.Command, _ []string) error {
	cfg, projectDir, err := project.Load()
	if err != nil {
		if errors.Is(err, project.ErrNotFound) {
			return errNoProject
		}

		return fmt.Errorf("load project: %w", err)
	}

	w := os.Stdout

	mustWrite(w, fmt.Sprintf("Project: %s\n\n", project.GetConfigPath(projectDir)))

	if len(cfg.Tools) > 0 {
		mustWrite(w, fmt.Sprintf("Tools: %s\n\n", strings.Join(cfg.Tools, ", ")))
	} else {
		mustWrite(w, "Tools: all compatible\n\n")
	}

	if len(cfg.Skills) > 0 {
		mustWrite(w, "Skills:\n")

		for _, s := range cfg.Skills {
			mustWrite(w, fmt.Sprintf("  - %s\n", s))
		}

		mustWrite(w, "\n")
	}

	if len(cfg.Agents) > 0 {
		mustWrite(w, "Agents:\n")

		for _, a := range cfg.Agents {
			mustWrite(w, fmt.Sprintf("  - %s\n", a))
		}

		mustWrite(w, "\n")
	}

	if len(cfg.Registries) > 0 {
		mustWrite(w, "Project registries:\n")

		for _, r := range cfg.Registries {
			if r.IsLocal() {
				mustWrite(w, fmt.Sprintf("  - %s: %s\n", r.Name, r.Path))
			} else if r.IsGit() {
				mustWrite(w, fmt.Sprintf("  - %s: %s\n", r.Name, r.URL))
			}
		}
	}

	if cfg.IsEmpty() {
		mustWrite(w, "No skills or agents defined.\n")
		mustWrite(w, "Use 'skillsmith project add <name>' to add items.\n")
	}

	return nil
}
