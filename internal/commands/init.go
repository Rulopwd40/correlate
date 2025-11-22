package commands

import (
	"log"

	"github.com/Rulopwd40/correlate/internal/core"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init [identifier] [project-type]",
	Short: "Initialize a new correlate project",
	Args:  cobra.ExactArgs(2),
	Long:  "Sets up a new correlate project in the current directory by creating necessary configuration files.",
	Run: func(cmd *cobra.Command, args []string) {
		runInit(cmd, args)
	},
}

func runInit(cmd *cobra.Command, args []string) {
	log.Println("Initializing correlate project...")

	orch, err := core.Get[*core.Orchestrator]()
	if err != nil {
		log.Println("Error getting orchestrator:", err)
	}

	library := args[1]
	identifier := args[0]

	err = orch.Init(library, identifier)
	if err != nil {
		return
	}
	log.Println("Correlate project initialized successfully.")
}
