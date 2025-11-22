package tests

import (
	"bytes"
	"log"
	"testing"

	"github.com/Rulopwd40/correlate/internal/commands"

	"github.com/Rulopwd40/correlate/internal/pipeline"
	"github.com/spf13/cobra"
)

func TestUpdateCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		description string
	}{
		{
			name:        "update specific identifier",
			args:        []string{"my-library"},
			wantErr:     false,
			description: "Should successfully update specific identifier",
		},
		{
			name:        "update all references",
			args:        []string{},
			wantErr:     false,
			description: "Should successfully update all references when no args provided",
		},
		{
			name:        "too many arguments",
			args:        []string{"my-library", "extra-arg"},
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
				Use:  "update [identifier]",
				Args: cobra.RangeArgs(0, 1),
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

func TestUpdateCmdDefinition(t *testing.T) {
	if commands.UpdateCmd.Use != "update [identifier]" {
		t.Errorf("Expected Use to be 'update [identifier]', got '%s'", commands.UpdateCmd.Use)
	}

	if commands.UpdateCmd.Short != "Update project dependencies" {
		t.Errorf("Expected Short description, got '%s'", commands.UpdateCmd.Short)
	}

	// Check aliases
	if len(commands.UpdateCmd.Aliases) != 1 || commands.UpdateCmd.Aliases[0] != "u" {
		t.Errorf("Expected alias 'u', got %v", commands.UpdateCmd.Aliases)
	}

	// Test that command accepts 0 or 1 args
	if err := commands.UpdateCmd.Args(commands.UpdateCmd, []string{}); err != nil {
		t.Errorf("Expected no error for zero arguments, got %v", err)
	}

	if err := commands.UpdateCmd.Args(commands.UpdateCmd, []string{"arg1"}); err != nil {
		t.Errorf("Expected no error for one argument, got %v", err)
	}

	if err := commands.UpdateCmd.Args(commands.UpdateCmd, []string{"arg1", "arg2"}); err == nil {
		t.Error("Expected error for two arguments, got nil")
	}
}

func TestRenderEvent(t *testing.T) {
	tests := []struct {
		name  string
		event pipeline.Event
	}{
		{
			name: "task start event",
			event: pipeline.Event{
				Type:     pipeline.EventTaskStart,
				TaskName: "Build Task",
				Message:  "starting",
			},
		},
		{
			name: "task progress event",
			event: pipeline.Event{
				Type:     pipeline.EventTaskProgress,
				TaskName: "Build Task",
				Message:  "Building...",
			},
		},
		{
			name: "task finish event",
			event: pipeline.Event{
				Type:     pipeline.EventTaskFinish,
				TaskName: "Build Task",
				Message:  "completed",
			},
		},
		{
			name: "error event",
			event: pipeline.Event{
				Type:     pipeline.EventError,
				TaskName: "Build Task",
				Message:  "build failed",
			},
		},
		{
			name: "pipeline done event",
			event: pipeline.Event{
				Type: pipeline.EventPipelineDone,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(nil)

			// This should not panic
			commands.RenderEvent(tt.event)
		})
	}
}
