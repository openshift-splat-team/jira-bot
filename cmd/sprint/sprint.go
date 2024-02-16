package sprint

import (
	"github.com/spf13/cobra"
)

var applyFlag bool

var cmdSprint = &cobra.Command{
	Use:   "sprint",
	Short: "Manage sprints",
	Long:  `This command allows you to manage sprints in your project management tool.`,
}

func Initialize(rootCmd *cobra.Command) {

	cmdMoveIssue.Flags().BoolVarP(&applyFlag, "apply", "a", false, "Apply the changes")
	cmdMoveInQuery.Flags().BoolVarP(&applyFlag, "apply", "a", false, "Apply the changes")
	cmdSprint.AddCommand(cmdMoveIssue)
	cmdSprint.AddCommand(cmdMoveInQuery)
	rootCmd.AddCommand(cmdSprint)
}
