package main

import (
	"fmt"
	"os"
	"strings"

	jira "github.com/andygrunwald/go-jira"
)

const (
	FieldStoryPoints   = "customfield_12310243"
	FieldStatusSummary = "customfield_12320841"
)

func getStoryPoints(totalMap map[string]interface{}) float64 {
	if points, exists := totalMap[FieldStoryPoints]; exists {
		if points != nil {
			return points.(float64)
		}
	}
	return 0
}

func main() {
	token := os.Getenv("JIRA_PERSONAL_ACCESS_TOKEN")
	if len(token) == 0 {
		fmt.Println("JIRA_PERSONAL_ACCESS_TOKEN must be exported")
		os.Exit(1)
	}
	tp := jira.BearerAuthTransport{
		Token: token,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), "https://issues.redhat.com/")
	issues, _, err := jiraClient.Issue.Search("filter = \"SPLAT Team - Epics 4.16\"", nil)
	if err != nil {
		fmt.Printf("unable to get epics: %v", err)
		os.Exit(1)
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
			points := getStoryPoints(childIssue.Fields.Unknowns)
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
			if getStoryPoints(unpointedIssue.Fields.Unknowns) == 0 {
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

		if statusSummary == issue.Fields.Unknowns[FieldStatusSummary] {
			fmt.Println("no update")
			continue
		}
		propertyMap := map[string]interface{}{
			"fields": map[string]interface{}{
				FieldStoryPoints:   aggregatePoints,
				FieldStatusSummary: statusSummary,
			},
		}
		resp, err := jiraClient.Issue.UpdateIssue(issue.Key, propertyMap)
		if err != nil {
			fmt.Printf("unable to update epic %s: %v\n", issue.Key, err)
			b := []byte{}
			resp.Response.Body.Read(b)
			fmt.Println(string(b))
			continue
		}
	}
}
