package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Rulopwd40/correlate/internal/core"
	"github.com/Rulopwd40/correlate/internal/pipeline"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:     "update [identifier]",
	Aliases: []string{"u"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   "Update project",
	Long:    "Concurrent process that gets references versions and update the project",
	Run: func(cmd *cobra.Command, args []string) {
		runUpdate(cmd, args)
	},
}

func runUpdate(cmd *cobra.Command, args []string) {
	log.Println("Updating project...")

	orch, err := core.Get[*core.Orchestrator]()
	if err != nil {
		log.Println("Error getting orchestrator:", err)
		return
	}

	var identifier string
	if len(args) > 0 {
		identifier = args[0]
	} else {
		identifier = "" // Update ALL references
	}

	err = orch.Update(identifier)
	if err != nil {
		log.Println("Error during update:", err)
		return
	}

	errorPresent := false
	for ev := range orch.Events() {
		if ev.Type == pipeline.EventError {
			errorPresent = true
		}
		RenderEvent(ev)
	}

	if errorPresent {
		log.Println("Error during update. Check Logs")
		return
	}
	log.Println("Project successfully updated!")
}
func RenderEvent(ev pipeline.Event) {
	switch ev.Type {

	case pipeline.EventTaskStart:
		fmt.Printf("START: %s\n", ev.TaskName)

	case pipeline.EventTaskProgress:
		line := strings.TrimSpace(ev.Message)
		if line != "" {
			fmt.Printf("   %s\n", line)
		}

	case pipeline.EventTaskFinish:
		fmt.Printf("DONE:  %s\n", ev.TaskName)

	case pipeline.EventError:
		fmt.Printf(" ERROR in %s: %s\n", ev.TaskName, ev.Message)
		if ev.Err != nil {
			fmt.Printf("   ↳ %v\n", ev.Err)
		}

	case pipeline.EventPipelineDone:
		fmt.Println("Pipeline finished successfully.")
	}
}
