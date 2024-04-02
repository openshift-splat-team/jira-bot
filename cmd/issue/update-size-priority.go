package issue

import (
	"fmt"
	"log"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	cmdUpdateSizeAndPriority.Flags().BoolVarP(&options.dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")
	cmdUpdateSizeAndPriority.Flags().BoolVarP(&options.overrideFlag, "override", "o", false, "overrides a warning when --override=true")
	cmdUpdateSizeAndPriority.Flags().Int64VarP(&options.points, "points", "p", -1, "points to apply to issue")
	cmdUpdateSizeAndPriority.Flags().StringVarP(&options.priority, "priority", "r", "", "priority to set")
	cmdUpdateSizeAndPriority.Flags().StringVarP(&options.comment, "comment", "c", "", "comment to append to issue")
	cmdUpdateSizeAndPriority.Flags().StringVarP(&options.resolution, "resolution", "", "", "resolution to set on an issue")
	cmdUpdateSizeAndPriority.Flags().StringVarP(&options.state, "state", "s", "", "sets the issue state")
	cmdIssue.AddCommand(cmdUpdateSizeAndPriority)
}

// addComment adds a comment to a jira card 
func addComment(client *jira.Client, issue jira.Issue, options *issueCommandOptions) error {
	//user, _, err := client.User.GetSelf() 
	//if err != nil {
	//	return fmt.Errorf("unable to get self: %v", err)
	//} 
	comment := &jira.Comment {
		Body: options.comment,
		//Author: *user,
	}

	if options.dryRunFlag {
		log.Printf("would have added comment to issue: %s", issue.Key)
	} else {
		log.Printf("adding comment to issue %s", issue.Key)
		_, _, err := client.Issue.AddComment(issue.Key, comment)
		return err
	}
	 
	return nil
}

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

	priorities := []string{}

	for _, knownPriority := range knownPriorities {
		priorities = append(priorities, knownPriority.Name)
		if strings.EqualFold(knownPriority.Name, strings.ToLower(options.priority)) {
			priority = &knownPriority
		}
	}

	if priority == nil {
		return fmt.Errorf("priority %s does not match a known priority: %s", options.priority, strings.Join(priorities, ","))
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

func checkSetState(client *jira.Client, issue jira.Issue, options *issueCommandOptions) error {
	log.Printf("attempting to transition issue to %s", options.state)
	transitions, _, err := client.Issue.GetTransitions(issue.ID)
	if err != nil {
		return fmt.Errorf("unable to get known transitions: %v", err)
	}

	var status *jira.Transition

	transitionList := []string{}
	for idx, transition := range transitions {
		transitionList = append(transitionList, transition.Name)
		if strings.EqualFold(transition.Name, strings.ToLower(options.state)) {
			status = &transitions[idx]
		}
	}

	if status == nil {
		return fmt.Errorf("state %s does not match a known transition: %s", options.state, strings.Join(transitionList, ","))
	}

	if options.dryRunFlag {
		log.Printf("issue: %s would have been updated. run again and provide --dry-run=false to apply.", issue.Key)
		return nil
	} else {
		log.Printf("transitioning to \"%s\" for issue: %s", status.Name, issue.Key)
		_, err := client.Issue.DoTransition(issue.Key, status.ID)
		if err != nil {
			return fmt.Errorf("unable to update issue %s: %v", issue.Key, err)
		}
	}

	return nil
}

func checkSetResolution(client *jira.Client, issue jira.Issue, options *issueCommandOptions) error {
	log.Printf("attempting to set resolution to %s", options.resolution)
	resolutions, _, err := client.Resolution.GetList()
	
	if err != nil {
		return fmt.Errorf("unable to get known transitions: %v", err)
	}

	var status *jira.Resolution  

	resolutionList := []string{}
	for idx, resolution := range resolutions {
		resolutionList = append(resolutionList, resolution.Name)
		if strings.EqualFold(resolution.Name, strings.ToLower(options.resolution)) {
			status = &resolutions[idx]
		}
	}

	if status == nil {
		return fmt.Errorf("resolution %s does not match a known resultion: %s", options.state, strings.Join(resolutionList, ","))
	}

	if options.dryRunFlag {
		log.Printf("issue: %s would have been updated. run again and provide --dry-run=false to apply.", issue.Key)
		return nil
	} else {
		propertyMap := map[string]interface{}{
			"fields": map[string]interface{}{
				"resolution": status,
			},
		}
		log.Printf("setting resolution to \"%s\" for issue: %s", status.Name, issue.Key)
		issue.Fields.Resolution = status
		client.Issue.UpdateIssue(issue.Key, propertyMap)
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
		if len(options.state) > 0 {
			err = checkSetState(jiraClient, issue, options)
			if err != nil {
				return fmt.Errorf("unable to set story state: %v", err)
			}
		}
		if len(options.resolution) > 0 {
			err = checkSetResolution(jiraClient, issue, options)
			if err != nil {
				return fmt.Errorf("unable to issue resolution: %v", err)
			}
		}
		if len(options.comment) > 0 {  
			err = addComment(jiraClient, issue, options)
			if err != nil {  
				return fmt.Errorf("unable to udpate story: %v", err)
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
