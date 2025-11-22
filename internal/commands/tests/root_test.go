package tests

import (
	"testing"

	"github.com/Rulopwd40/correlate/internal/commands"
)

func TestRootCmd(t *testing.T) {
	if commands.RootCmd.Use != "correlate" {
		t.Errorf("Expected Use to be 'correlate', got '%s'", commands.RootCmd.Use)
	}

	if commands.RootCmd.Short != "Correlate - Dependency management automation tool" {
		t.Errorf("Expected Short description, got '%s'", commands.RootCmd.Short)
	}
}

func TestRootCommandHasSubcommands(t *testing.T) {
	commands := commands.RootCmd.Commands()

	if len(commands) == 0 {
		t.Error("Expected root command to have subcommands, got none")
	}

	expectedCommands := map[string]bool{
		"init":    false,
		"link":    false,
		"replace": false,
		"update":  false,
	}

	for _, cmd := range commands {
		if _, exists := expectedCommands[cmd.Name()]; exists {
			expectedCommands[cmd.Name()] = true
		}
	}

	for cmdName, found := range expectedCommands {
		if !found {
			t.Errorf("Expected subcommand '%s' not found", cmdName)
		}
	}
}

func TestRootCommandStructure(t *testing.T) {
	// Test that Execute function exists and doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked: %v", r)
		}
	}()

	// We can't really execute the command in tests without mocking
	// but we can verify the structure exists
	if commands.RootCmd == nil {
		t.Error("rootCmd should not be nil")
	}
}
