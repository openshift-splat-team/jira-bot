package sprint

import (
	"github.com/spf13/cobra"
)

var dryRunFlag bool
var overrideFlag bool

var cmdSprint = &cobra.Command{
	Use:   "sprint",
	Short: "Manage sprints",
	Long:  `This command allows you to manage sprints in your project management tool.`,
}

func Initialize(rootCmd *cobra.Command) {

	cmdMoveIssue.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")
	cmdMoveInQuery.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")
	cmdMoveIssue.Flags().BoolVarP(&overrideFlag, "override", "o", false, "adds issue to a sprint regardless if that issue is already in a sprint with --override=true")
	cmdMoveInQuery.Flags().BoolVarP(&overrideFlag, "override", "o", false, "adds issue to a sprint regardless if that issue is already in a sprint with --override=true")

	cmdSprint.AddCommand(cmdMoveIssue)
	cmdSprint.AddCommand(cmdMoveInQuery)
	rootCmd.AddCommand(cmdSprint)
}
