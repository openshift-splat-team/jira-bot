package issue

import (
	"fmt"
	"log"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/splat-jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

func checkSetPoints(client *jira.Client, issue jira.Issue, options *issueCommandOptions) error {
	if util.GetStoryPoints(issue.Fields.Unknowns) > 0 && !options.overrideFlag {
		log.Fatalf("issue: %s already has assigned story points.  run again and provide --override=true to apply", issue.Key)
		return nil
	}

	if util.GetStoryPoints(issue.Fields.Unknowns) == 0 || options.overrideFlag {
		propertyMap := map[string]interface{}{
			"fields": map[string]interface{}{
				util.FieldStoryPoints: options.points,
			},
		}

		if options.dryRunFlag {
			log.Printf("issue: %s would have been updated. run again and provide --dry-run=false to apply.", issue.Key)
			return nil
		} else {
			log.Printf("setting story points for issue: %s", issue.Key)
			_, err := client.Issue.UpdateIssue(issue.Key, propertyMap)
			if err != nil {
				return fmt.Errorf("unable to update issue %s: %v", issue.Key, err)
			}
		}
	}

	return nil
}

func checkSetPriority(client *jira.Client, issue jira.Issue, options *issueCommandOptions) error {
	log.Printf("attempting to set issue priority")
	if issue.Fields.Priority != nil && !options.overrideFlag {
		log.Fatalf("issue: %s already has assigned priority.  run again and provide --override=true to apply", issue.Key)
		return nil
	}
	knownPriorities, _, err := client.Priority.GetList()
	if err != nil {
		return fmt.Errorf("unable to get known priorities: %v", err)
	}

	var priority *jira.Priority

	for _, knownPriority := range knownPriorities {
		if strings.EqualFold(knownPriority.Name, strings.ToLower(options.priority)) {
			priority = &knownPriority
			break
		}
	}

	if priority == nil {
		return fmt.Errorf("priority %s does not match a known priority", options.priority)
	}

	if issue.Fields.Priority == nil || options.overrideFlag {
		propertyMap := map[string]interface{}{
			"fields": map[string]interface{}{
				"priority": priority,
			},
		}

		if options.dryRunFlag {
			log.Printf("issue: %s would have been updated. run again and provide --dry-run=false to apply.", issue.Key)
			return nil
		} else {
			log.Printf("setting priority to %s for issue: %s", options.priority, issue.Key)
			_, err := client.Issue.UpdateIssue(issue.Key, propertyMap)
			if err != nil {
				return fmt.Errorf("unable to update issue %s: %v", issue.Key, err)
			}
		}
	}

	return nil
}

// updateSizeAndPriority according to rules set forth by the team
func updateSizeAndPriority(issue string, options *issueCommandOptions) error {
	log.Printf("preparing to update issue: %s", issue)
	jiraClient, err := util.GetJiraClient()
	if err != nil {
		return fmt.Errorf("unable to get Jira client: %v", err)
	}

	issues, _, err := jiraClient.Issue.Search(fmt.Sprintf("issuekey in (%s)", issue), nil)
	if err != nil {
		return fmt.Errorf("unable to get issues: %v", err)
	}

	log.Printf("%d issues found in query", len(issues))

	for _, issue := range issues {
		if options.points != -1 {
			err = checkSetPoints(jiraClient, issue, options)
			if err != nil {
				return fmt.Errorf("unable to set story points: %v", err)
			}
		}
		if len(options.priority) > 0 {
			err = checkSetPriority(jiraClient, issue, options)
			if err != nil {
				return fmt.Errorf("unable to set story priority: %v", err)
			}
		}
	}
	return nil
}

var cmdUpdateSizeAndPriority = &cobra.Command{
	Use:   "update-size-and-priority [issue]",
	Short: "Updates issues according to rules provided as options.",
	Long:  `Updates issues matching the JQL provided according to rules provided as options`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := updateSizeAndPriority(args[0], &options)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to update issue: %v", err))
		}
	},
}
