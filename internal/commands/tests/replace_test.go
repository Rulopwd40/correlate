package tests

import (
	"bytes"
	"github.com/Rulopwd40/correlate/internal/commands"
	"log"
	"testing"

	"github.com/spf13/cobra"
)

func TestReplaceCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		wantErr     bool
		description string
	}{
		{
			name:        "valid replace command without version",
			args:        []string{"my-library"},
			flags:       map[string]string{},
			wantErr:     false,
			description: "Should successfully replace with just identifier",
		},
		{
			name:        "valid replace command with version flag",
			args:        []string{"my-library"},
			flags:       map[string]string{"version": "1.2.3"},
			wantErr:     false,
			description: "Should successfully replace with identifier and version",
		},
		{
			name:        "no arguments",
			args:        []string{},
			flags:       map[string]string{},
			wantErr:     true,
			description: "Should fail when no arguments provided",
		},
		{
			name:        "too many arguments",
			args:        []string{"my-library", "extra-arg"},
			flags:       map[string]string{},
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
				Use:  "replace [identifier]",
				Args: cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().StringP("version", "v", "", "Version to apply (optional)")

			// Set flags
			for key, value := range tt.flags {
				cmd.Flags().Set(key, value)
			}

			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("%s: error = %v, wantErr %v", tt.description, err, tt.wantErr)
			}
		})
	}
}

func TestReplaceCmdDefinition(t *testing.T) {
	if commands.ReplaceCmd.Use != "replace [identifier]" {
		t.Errorf("Expected Use to be 'replace [identifier]', got '%s'", commands.ReplaceCmd.Use)
	}

	if commands.ReplaceCmd.Short != "Replace version in project dependency" {
		t.Errorf("Expected Short description, got '%s'", commands.ReplaceCmd.Short)
	}

	// Test that command accepts exactly 1 arg
	if err := commands.ReplaceCmd.Args(commands.ReplaceCmd, []string{}); err == nil {
		t.Error("Expected error for no arguments, got nil")
	}

	if err := commands.ReplaceCmd.Args(commands.ReplaceCmd, []string{"arg1"}); err != nil {
		t.Errorf("Expected no error for one argument, got %v", err)
	}

	if err := commands.ReplaceCmd.Args(commands.ReplaceCmd, []string{"arg1", "arg2"}); err == nil {
		t.Error("Expected error for two arguments, got nil")
	}

	// Check version flag exists
	versionFlag := commands.ReplaceCmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Error("Expected 'version' flag to exist")
	}
	if versionFlag.Shorthand != "v" {
		t.Errorf("Expected shorthand 'v', got '%s'", versionFlag.Shorthand)
	}
}
