package issue

import (
	"fmt"
	"log"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/splat-jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

const (
	FieldStoryPoints   = "customfield_12310243"
	FieldStatusSummary = "customfield_12320841"
)

type autoUpdateIssueOptions struct {
	defaultSpikeStoryPoints int64
	dryRunFlag              bool
	overrideFlag            bool
}

func (a *autoUpdateIssueOptions) checkSetSpikePoints(client *jira.Client, issue jira.Issue) error {
	if issue.Fields.Type.Name == "Spike" {
		if util.GetStoryPoints(issue.Fields.Unknowns) > 0 && !a.overrideFlag {
			log.Fatalf("issue: %s already has assigned story points.  run again and provide --override=true to apply", issue.Key)
			return nil
		}

		if util.GetStoryPoints(issue.Fields.Unknowns) == 0 || a.overrideFlag {
			propertyMap := map[string]interface{}{
				"fields": map[string]interface{}{
					util.FieldStoryPoints: a.defaultSpikeStoryPoints,
				},
			}
			if a.dryRunFlag {
				log.Printf("issue: %s would have default spike points assigned. run again and provide --dry-run=false to apply.", issue.Key)
				return nil
			} else {
				log.Printf("setting default story points for spike: %s", issue.Key)
				_, err := client.Issue.UpdateIssue(issue.Key, propertyMap)
				if err != nil {
					return fmt.Errorf("unable to update issue %s: %v", issue.Key, err)
				}
			}
		}
	}
	return nil
}

// autoUpdateIssuesInQuery according to rules set forth by the team
func (a *autoUpdateIssueOptions) autoUpdateIssuesInQuery(jql string) error {
	log.Printf("preparing to auto-update issues found in query: %s", jql)
	jiraClient, err := util.GetJiraClient()
	if err != nil {
		return fmt.Errorf("unable to get Jira client: %v", err)
	}

	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		return fmt.Errorf("unable to get issues: %v", err)
	}

	log.Printf("%d issues found in query", len(issues))

	for _, issue := range issues {
		if a.defaultSpikeStoryPoints > 0 {
			err = a.checkSetSpikePoints(jiraClient, issue)
			if err != nil {
				return fmt.Errorf("unable to set default spike story points: %v", err)
			}
		}
	}
	return nil
}

func Initialize(rootCmd *cobra.Command) {
	options := autoUpdateIssueOptions{}

	cmdAutoUpdateIssuesStatus := &cobra.Command{
		Use:   "auto-update-issues [jql]",
		Short: "Updates issues according to rules provided as options.",
		Long:  `Updates issues matching the JQL provided according to rules provided as options`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := options.autoUpdateIssuesInQuery(args[0])
			if err != nil {
				util.RuntimeError(fmt.Errorf("unable to update issues: %v", err))
			}
		},
	}
	cmdAutoUpdateIssuesStatus.Flags().BoolVarP(&options.dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")
	cmdAutoUpdateIssuesStatus.Flags().BoolVarP(&options.overrideFlag, "override", "o", false, "overrides a warning when --override=true")
	cmdAutoUpdateIssuesStatus.Flags().Int64VarP(&options.defaultSpikeStoryPoints, "default-spike-points", "s", -1, "points to apply to spikes which have no points")
	rootCmd.AddCommand(cmdAutoUpdateIssuesStatus)
}
