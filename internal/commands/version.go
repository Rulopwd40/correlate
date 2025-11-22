package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Display the version, commit hash, and build date of correlate.

This information is useful for debugging and ensuring you're running
the expected version of the tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("correlate version %s\n", Version)
		fmt.Printf("  commit: %s\n", Commit)
		fmt.Printf("  built:  %s\n", BuildDate)
	},
}
