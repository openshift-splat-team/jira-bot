package issue

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strconv"
)

func init() {
	cmdIssue.AddCommand(cmdGenerateIssueSizings)
}

var re = regexp.MustCompile(`[\r\n,\t]`)

func generateIssueSizings(filter string, options *issueCommandOptions) error {
	client, err := util.GetJiraClient()
	if err != nil {
		return fmt.Errorf("unable to get Jira client: %v", err)
	}

	issues, _, err := client.Issue.Search(filter, &jira.SearchOptions{})
	if err != nil {
		return fmt.Errorf("unable to search for issues: %v", err)
	}

	csvDict := map[string][]string{"Description": make([]string, 0), "Summary": make([]string, 0), "Issue Type": make([]string, 0), "Issue Key": make([]string, 0)}
	for _, issue := range issues {
		if issue.Fields.Type.Name == "Bug" {
			fmt.Printf("Issue %s is a bug, bugs don't get sizings.  skipping...\n", issue.Key)
			continue
		}

		if util.GetStoryPoints(issue.Fields.Unknowns) > 0 {
			fmt.Printf("Issue %s already has a sizing.  skipping...\n", issue.Key)
			continue
		}

		description := issue.Fields.Description
		description = re.ReplaceAllString(description, " ")
		summary := issue.Fields.Summary
		summary = re.ReplaceAllString(summary, " ")

		csvDict["Description"] = append(csvDict["Description"], description)
		csvDict["Summary"] = append(csvDict["Summary"], summary)
		csvDict["Issue Type"] = append(csvDict["Issue Type"], issue.Fields.Type.Name)
		csvDict["Issue Key"] = append(csvDict["Issue Key"], issue.Key)
	}

	url := os.Getenv("JIRA_NEURAL_SIZING_URL")
	if len(url) == 0 {
		url = "http://127.0.0.1:8001"
	}
	outMap, err := util.PostJSONData(url, csvDict)
	if err != nil {
		return fmt.Errorf("error sending request to %s: %v", url, err)
	}

	for _, issue := range issues {
		foundIssue := false
		for k, v := range outMap["Issue Key"] {
			if v != issue.Key {
				continue
			}
			val, err := strconv.Atoi(outMap["sizing"][k])
			if err != nil {
				fmt.Printf("Issue %s is not a valid size.  skipping...\n", issue.Key)
				continue
			}
			foundIssue = true
			options.points = int64(val)
			break
		}
		if !foundIssue {
			continue
		}

		updated, err := checkSetPoints(client, issue, options)
		if err != nil {
			return fmt.Errorf("error setting points on issue: %v", err)
		}

		if updated {
			_, _, err = client.Issue.AddComment(issue.Key, &jira.Comment{
				Body: "set sizing based on past history of issues sized by the team. this issue should still be sized and refined.",
				Visibility: jira.CommentVisibility{
					Type:  "role",
					Value: "Administrators",
				}})
			if err != nil {
				return fmt.Errorf("unable to add comment %s: %v", url, err)
			}
		}
		break
	}
	return nil
}

var cmdGenerateIssueSizings = &cobra.Command{
	Use:   "export-as-csv [filter]",
	Short: "Exports a query as CSV",
	Long:  `Exports a query as CSV`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := generateIssueSizings(args[0], &options)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to perform query: %v", err))
		}
	},
}
