package util

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

// Custom field IDs for on-premise Jira
const (
	OnPremFieldStoryPoints   = "customfield_12310243"
	OnPremFieldStatusSummary = "customfield_12320841"
	OnPremFieldEpicLink      = "customfield_12311140"
	OnPremFieldFeatureLink   = "customfield_12318341"
)

// Custom field IDs for Atlassian Cloud
// Note: These may vary by organization - verify with your Jira instance
const (
	CloudFieldStoryPoints   = "customfield_10028"
	CloudFieldStatusSummary = "customfield_10841"
	CloudFieldEpicLink      = "customfield_10014"
	CloudFieldFeatureLink   = "customfield_10341"
)

// Dynamic field accessors
var (
	FieldStoryPoints   string
	FieldStatusSummary string
	FieldEpicLink      string
	FieldFeatureLink   string
)

func isCloudInstance(baseURL string) bool {
	return strings.Contains(strings.ToLower(baseURL), ".atlassian.net")
}

func initializeFieldIDs(baseURL string) {
	if isCloudInstance(baseURL) {
		FieldStoryPoints = CloudFieldStoryPoints
		FieldStatusSummary = CloudFieldStatusSummary
		FieldEpicLink = CloudFieldEpicLink
		FieldFeatureLink = CloudFieldFeatureLink
		log.Println("Using Atlassian Cloud field IDs")
	} else {
		FieldStoryPoints = OnPremFieldStoryPoints
		FieldStatusSummary = OnPremFieldStatusSummary
		FieldEpicLink = OnPremFieldEpicLink
		FieldFeatureLink = OnPremFieldFeatureLink
		log.Println("Using on-premise Jira field IDs")
	}
}

func GetJiraClient() (*jira.Client, error) {
	// Get base URL from environment or use default
	baseURL := viper.GetString("base_url")
	if len(baseURL) == 0 {
		baseURL = "https://issues.redhat.com"
		log.Printf("JIRA_BASE_URL not set, using default: %s\n", baseURL)
	}

	// Ensure base URL ends with /
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// Initialize custom field IDs based on instance type
	initializeFieldIDs(baseURL)

	var client *jira.Client
	var err error

	if isCloudInstance(baseURL) {
		// Atlassian Cloud: Use BasicAuth with email + API token
		email := viper.GetString("email")
		apiToken := viper.GetString("api_token")

		if len(email) == 0 || len(apiToken) == 0 {
			return nil, fmt.Errorf("Atlassian Cloud requires JIRA_EMAIL and JIRA_API_TOKEN environment variables")
		}

		log.Printf("Connecting to Atlassian Cloud: %s (user: %s)\n", baseURL, email)

		tp := jira.BasicAuthTransport{
			Username: email,
			APIToken: apiToken,
		}

		client, err = jira.NewClient(tp.Client(), baseURL)
	} else {
		// On-premise: Use Bearer token
		token := viper.GetString("personal_access_token")

		if len(token) == 0 {
			return nil, fmt.Errorf("On-premise Jira requires JIRA_PERSONAL_ACCESS_TOKEN environment variable")
		}

		log.Printf("Connecting to on-premise Jira: %s\n", baseURL)

		tp := jira.BearerAuthTransport{
			Token: token,
		}

		client, err = jira.NewClient(tp.Client(), baseURL)
	}

	return client, err
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

// GetParentLinks retrieves the parent links from the total map
func GetParentLinks(totalMap map[string]interface{}) (string, string) {	
	var epicLink, featureLink string
	var val interface{}
	val = totalMap[FieldEpicLink]
	if val != nil {
		epicLink = val.(string)
	}
	val = totalMap[FieldFeatureLink]
	if val != nil {
		featureLink = val.(map[string]interface{})["key"].(string)
	}
	log.Printf("epic link: %s -- feature link: %s", epicLink, featureLink)
	return epicLink, featureLink
}


// GetIssueType retrieves the identified issue type from Jira
func GetIssueType(project *jira.Project, typeID string) (*jira.IssueType, error) {
	log.Printf("getting issue type: %s", typeID)

	for _, issueType := range project.IssueTypes {
		if strings.EqualFold(issueType.Name, typeID) {
			return &issueType, nil
		}
	}

	return nil, fmt.Errorf("unable to find issue type: %s", typeID)
}

// GetUser retrieves the identified user from Jira
func GetUser(client *jira.Client, userID string) (*jira.User, error) {
	log.Printf("getting user: %s", userID)

	user, _, err := client.User.GetSelf()
	if err != nil {
		return nil, fmt.Errorf("unable to get user: %v", err)
	}
	return user, nil
}

// GetProject retrieves the identified project from Jira
func GetProject(client *jira.Client, projectID string) (*jira.Project, error) {
	log.Printf("getting project: %s", projectID)
	project, _, err := client.Project.Get(projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to get project: %v", err)
	}
	return project, nil
}

func GetResponseBody(resp *jira.Response) (string, error) {
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %v", err)

	}
	return string(body), nil
}

func GetUnsizedStories() ([]jira.Issue, error) {
	client, err := GetJiraClient()
	if err != nil {
		return nil, fmt.Errorf("unable to get Jira client: %v", err)
	}
	unpointedIssues, _, err := client.Issue.Search("filter = \"OpenShift SPLAT - No story points assigned\"", nil)
	return unpointedIssues, err
}