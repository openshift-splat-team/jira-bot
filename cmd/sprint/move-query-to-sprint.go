package sprint

import (
	"fmt"

	"github.com/openshift-splat-team/splat-jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

var cmdMoveInQuery = &cobra.Command{
	Use:   "move-in-query [sprint-number] [JQL query]...",
	Short: "Moves all issues returned by a JQL query",
	Long:  `This command allows you to remove an issue from a sprint in your project management tool.`,
	Args:  cobra.MinimumNArgs(2), // Requires exactly two arguments: sprint-number and issue-number
	Run: func(cmd *cobra.Command, args []string) {
		err := util.CheckForMissingEnvVars()
		if err != nil {
			util.RuntimeError(err)
		}
		sprintNumber := args[0]
		query := args[1]
		client, err := util.GetJiraClient()
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to get jira client: %v", err))
		}
		_, issueIds, err := util.GetIssuesInQuery(client, query)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to get issues in query: %v", err))
		}

		err = moveToSprint(client, sprintNumber, issueIds)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to move issues to sprint: %v", err))
		}
	},
}
