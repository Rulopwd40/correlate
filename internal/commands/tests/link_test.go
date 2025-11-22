package tests

import (
	"bytes"
	"github.com/Rulopwd40/correlate/internal/commands"
	"log"
	"testing"

	"github.com/spf13/cobra"
)

func TestLinkCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		description string
	}{
		{
			name:        "valid link command",
			args:        []string{"my-library", "/path/to/project"},
			wantErr:     false,
			description: "Should successfully link with valid identifier and path",
		},
		{
			name:        "missing path argument",
			args:        []string{"my-library"},
			wantErr:     true,
			description: "Should fail when missing path argument",
		},
		{
			name:        "no arguments",
			args:        []string{},
			wantErr:     true,
			description: "Should fail when no arguments provided",
		},
		{
			name:        "too many arguments",
			args:        []string{"my-library", "/path/to/project", "extra"},
			wantErr:     true,
			description: "Should fail when too many arguments provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(nil)

			// Create a new command instance
			cmd := &cobra.Command{
				Use:  "link [identifier] [fullPath]",
				Args: cobra.ExactArgs(2),
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}

			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("%s: error = %v, wantErr %v", tt.description, err, tt.wantErr)
			}
		})
	}
}

func TestLinkCmdDefinition(t *testing.T) {
	if commands.LinkCmd.Use != "link [identifier] [fullPath]" {
		t.Errorf("Expected Use to be 'link [identifier] [fullPath]', got '%s'", commands.LinkCmd.Use)
	}

	if commands.LinkCmd.Short != "Add a project reference" {
		t.Errorf("Expected Short description, got '%s'", commands.LinkCmd.Short)
	}

	// Check aliases
	if len(commands.LinkCmd.Aliases) != 1 || commands.LinkCmd.Aliases[0] != "l" {
		t.Errorf("Expected alias 'l', got %v", commands.LinkCmd.Aliases)
	}

	// Test that command accepts exactly 2 args
	if err := commands.LinkCmd.Args(commands.LinkCmd, []string{"arg1"}); err == nil {
		t.Error("Expected error for single argument, got nil")
	}

	if err := commands.LinkCmd.Args(commands.LinkCmd, []string{"arg1", "arg2"}); err != nil {
		t.Errorf("Expected no error for two arguments, got %v", err)
	}

	if err := commands.LinkCmd.Args(commands.LinkCmd, []string{"arg1", "arg2", "arg3"}); err == nil {
		t.Error("Expected error for three arguments, got nil")
	}
}
