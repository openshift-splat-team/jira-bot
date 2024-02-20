package util

import (
	"fmt"
	"log"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

const (
	FieldStoryPoints   = "customfield_12310243"
	FieldStatusSummary = "customfield_12320841"
)

func GetJiraClient() (*jira.Client, error) {
	token := viper.GetString("personal_access_token")

	tp := jira.BearerAuthTransport{
		Token: token,
	}

	return jira.NewClient(tp.Client(), "https://issues.redhat.com/")
}

func GetIssuesInQuery(client *jira.Client, query string) ([]jira.Issue, []string, error) {
	log.Printf("invoking query: %s\n", query)
	issues, _, err := client.Issue.Search(query, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to execute query: %v", err)
	}
	issueIds := []string{}
	for _, issue := range issues {
		issueIds = append(issueIds, issue.ID)
	}
	log.Printf("found %d issues\n", len(issues))
	return issues, issueIds, nil
}

func GetStoryPoints(totalMap map[string]interface{}) float64 {
	if points, exists := totalMap[FieldStoryPoints]; exists {
		if points != nil {
			return points.(float64)
		}
	}
	return 0
}
