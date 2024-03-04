package epic

import (
	"fmt"
	"strings"

	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

func updateEpicStatus() error {
	jiraClient, err := util.GetJiraClient()
	if err != nil {
		return fmt.Errorf("unable to get Jira client: %v", err)
	}

	issues, _, err := jiraClient.Issue.Search("filter = \"SPLAT Team - Epics 4.16\"", nil)
	if err != nil {
		return fmt.Errorf("unable to get epics: %v", err)
	}

	for _, issue := range issues {
		fmt.Printf("\nchecking epic %s\n", issue.Fields.Summary)
		childIssues, _, err := jiraClient.Issue.Search(fmt.Sprintf("\"Parent Link\" = \"%s\"", issue.Key), nil)
		if err != nil {
			fmt.Printf("unable to get child issues for %s: %v", issue.Key, err)
			continue
		}

		aggregatePoints := 0.0
		completedPoints := 0.0
		inprogressPoints := 0.0
		unsizedStories := 0

		for _, childIssue := range childIssues {
			points := util.GetStoryPoints(childIssue.Fields.Unknowns)
			aggregatePoints += points
			if childIssue.Fields.Status.Name == "Closed" {
				completedPoints += points
			} else if childIssue.Fields.Status.Name != "Backlog" {
				inprogressPoints += points
			}
		}

		unpointedIssues, _, err := jiraClient.Issue.Search(fmt.Sprintf("filter = \"OpenShift SPLAT - No story points assigned\" AND \"Parent Link\" = \"%s\"", issue.Key), nil)
		if err != nil {
			fmt.Printf("unable to get child issues for %s: %v", issue.Key, err)
			continue
		}

		for _, unpointedIssue := range unpointedIssues {
			if util.GetStoryPoints(unpointedIssue.Fields.Unknowns) == 0 {
				unsizedStories++
			}
		}

		messages := []string{}

		if unsizedStories > 0 {
			messages = append(messages, fmt.Sprintf("unsized stories: %d", unsizedStories))
		}

		if completedPoints == aggregatePoints && unsizedStories == 0 {
			messages = []string{}
		} else {
			messages = append(messages, fmt.Sprintf("C/I/T: %.0f/%.0f/%.0f", completedPoints, inprogressPoints, aggregatePoints))
		}

		statusSummary := strings.Join(messages, "\n")

		if statusSummary == issue.Fields.Unknowns[util.FieldStatusSummary] {
			fmt.Println("no update")
			continue
		}
		propertyMap := map[string]interface{}{
			"fields": map[string]interface{}{
				util.FieldStoryPoints:   aggregatePoints,
				util.FieldStatusSummary: statusSummary,
			},
		}
		_, err = jiraClient.Issue.UpdateIssue(issue.Key, propertyMap)
		if err != nil {
			fmt.Printf("unable to update epic %s: %v\n", issue.Key, err)
			continue
		}
	}
	return nil
}

var cmdUpdateEpicStatus = &cobra.Command{
	Use:   "update-epic-status",
	Short: "Update the status of an epic",
	Long:  `This command automates the generation of epic status.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := updateEpicStatus()
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to update epic status: %v", err))
		}
	},
}

func Initialize(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cmdUpdateEpicStatus)
}
