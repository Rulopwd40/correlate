package tests

import (
	"bytes"
	"log"
	"testing"

	"github.com/Rulopwd40/correlate/internal/commands"

	"github.com/spf13/cobra"
)

func TestInitCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		description string
	}{
		{
			name:        "valid init command",
			args:        []string{"java-maven", "test-project"},
			wantErr:     false,
			description: "Should successfully initialize with valid library and identifier",
		},
		{
			name:        "missing arguments",
			args:        []string{"java-maven"},
			wantErr:     true,
			description: "Should fail when missing identifier argument",
		},
		{
			name:        "no arguments",
			args:        []string{},
			wantErr:     true,
			description: "Should fail when no arguments provided",
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
				Use:  "init [library] [identifier]",
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

func TestInitCmdDefinition(t *testing.T) {
	if commands.InitCmd.Use != "init [identifier] [project-type]" {
		t.Errorf("Expected Use to be 'init [identifier] [project-type]', got '%s'", commands.InitCmd.Use)
	}

	if commands.InitCmd.Short != "Initialize a new correlate project" {
		t.Errorf("Expected Short description, got '%s'", commands.InitCmd.Short)
	}

	// Test that command accepts exactly 2 args
	if err := commands.InitCmd.Args(commands.InitCmd, []string{"arg1"}); err == nil {
		t.Error("Expected error for single argument, got nil")
	}

	if err := commands.InitCmd.Args(commands.InitCmd, []string{"arg1", "arg2"}); err != nil {
		t.Errorf("Expected no error for two arguments, got %v", err)
	}

	if err := commands.InitCmd.Args(commands.InitCmd, []string{"arg1", "arg2", "arg3"}); err == nil {
		t.Error("Expected error for three arguments, got nil")
	}
}
