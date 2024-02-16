package sprint

import (
	"github.com/spf13/cobra"
)

var dryRunFlag bool

var cmdSprint = &cobra.Command{
	Use:   "sprint",
	Short: "Manage sprints",
	Long:  `This command allows you to manage sprints in your project management tool.`,
}

func Initialize(rootCmd *cobra.Command) {

	cmdMoveIssue.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")
	cmdMoveInQuery.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")
	cmdSprint.AddCommand(cmdMoveIssue)
	cmdSprint.AddCommand(cmdMoveInQuery)
	rootCmd.AddCommand(cmdSprint)
}
