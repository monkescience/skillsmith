package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

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

func setupCommands() {
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(listCmd)
}

//nolint:gochecknoinits // cobra requires init for command setup
func init() {
	setupCommands()
}

func runTUI(_ *cobra.Command, _ []string) error {
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("load registry: %w", err)
	}

	model := tui.NewModel(reg)
	p := tea.NewProgram(model, tea.WithAltScreen())

	_, err = p.Run()
	if err != nil {
		return fmt.Errorf("run tui: %w", err)
	}

	return nil
}

func runList(_ *cobra.Command, _ []string) error {
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("load registry: %w", err)
	}

	w := os.Stdout

	writeListOutput(w, reg)

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
