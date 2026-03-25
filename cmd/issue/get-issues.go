package issue

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
)

func init() {
	cmdGetIssues.Flags().BoolVarP(&jsonOutput, "json", "j", false, "dump JSON representation of issue(s)")
	cmdIssue.AddCommand(cmdGetIssues)
}

func getIssues(query string) error {
	client, err := util.GetJiraClient()
	if err != nil {
		return fmt.Errorf("unable to get Jira client: %v", err)
	}

	// Check if query is a single issue key or a JQL query
	var issues []jira.Issue
	if isSingleIssueKey(query) {
		// Fetch single issue by key
		issue, _, err := client.Issue.Get(query, nil)
		if err != nil {
			return fmt.Errorf("unable to get issue %s: %v", query, err)
		}
		issues = []jira.Issue{*issue}
	} else {
		// Execute JQL query
		searchIssues, _, err := client.Issue.Search(query, nil)
		if err != nil {
			return fmt.Errorf("unable to search for issues: %v", err)
		}
		issues = searchIssues
	}

	if len(issues) == 0 {
		log.Println("No issues found")
		return nil
	}

	if jsonOutput {
		// Dump JSON representation
		jsonBytes, err := json.MarshalIndent(issues, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal issues to JSON: %v", err)
		}
		fmt.Println(string(jsonBytes))
	} else {
		// Display formatted output
		for _, issue := range issues {
			fmt.Printf("Key: %s\n", issue.Key)
			fmt.Printf("Summary: %s\n", issue.Fields.Summary)
			fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
			fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
			if issue.Fields.Assignee != nil {
				fmt.Printf("Assignee: %s\n", issue.Fields.Assignee.DisplayName)
			} else {
				fmt.Printf("Assignee: Unassigned\n")
			}
			if issue.Fields.Priority != nil {
				fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
			}
			points := util.GetStoryPoints(issue.Fields.Unknowns)
			if points > 0 {
				fmt.Printf("Story Points: %.0f\n", points)
			}
			baseURL := client.GetBaseURL()
			fmt.Printf("URL: %sbrowse/%s\n", baseURL.String(), issue.Key)
			fmt.Println("---")
		}
		log.Printf("Found %d issue(s)", len(issues))
	}

	return nil
}

// isSingleIssueKey checks if the query looks like a single issue key (e.g., PROJ-123)
// rather than a JQL query
func isSingleIssueKey(query string) bool {
	// Simple heuristic: if it contains spaces or special JQL keywords, it's a query
	// Otherwise treat it as an issue key
	jqlKeywords := []string{"=", "AND", "OR", "NOT", "IN", "IS", "WAS", "ORDER BY", "filter"}
	for _, keyword := range jqlKeywords {
		if containsIgnoreCase(query, keyword) {
			return false
		}
	}
	return true
}

func containsIgnoreCase(s, substr string) bool {
	// Simple case-insensitive contains check
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

var cmdGetIssues = &cobra.Command{
	Use:   "get-issues [issue-key or JQL query]",
	Short: "Get issue(s) by key or JQL query",
	Long:  `Fetches and displays Jira issue(s) by issue key or JQL query. Use --json to dump the JSON representation.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := getIssues(args[0])
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to get issues: %v", err))
		}
	},
}
