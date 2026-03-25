package util

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var bindEnvVars = []string{"personal_access_token", "email", "api_token", "base_url", "project", "board"}

func CheckForMissingEnvVars() error {
	baseURL := viper.GetString("base_url")
	if len(baseURL) == 0 {
		baseURL = "https://issues.redhat.com"
	}

	// Check if this is an Atlassian Cloud instance
	isCloud := strings.Contains(strings.ToLower(baseURL), ".atlassian.net")

	// Validate authentication based on instance type
	if isCloud {
		// Cloud requires email and API token
		if len(viper.GetString("email")) == 0 || len(viper.GetString("api_token")) == 0 {
			return fmt.Errorf("Atlassian Cloud requires JIRA_EMAIL and JIRA_API_TOKEN environment variables")
		}
	} else {
		// On-prem requires personal access token
		if len(viper.GetString("personal_access_token")) == 0 {
			return fmt.Errorf("On-premise Jira requires JIRA_PERSONAL_ACCESS_TOKEN environment variable")
		}
	}

	// Check required variables
	requiredVars := []string{"project", "board"}
	for _, envVar := range requiredVars {
		if len(viper.GetString(envVar)) == 0 {
			return fmt.Errorf("the environment variable: [%s] must be exported", strings.ToUpper(fmt.Sprintf("jira_%s", envVar)))
		}
	}
	return nil
}

func BindEnvVars() error {
	viper.SetEnvPrefix("jira") // Set a prefix for environment variables
	for _, envVar := range bindEnvVars {
		err := viper.BindEnv(envVar)
		if err != nil {
			return fmt.Errorf("unable to bind env var %s: %v", envVar, err)
		}
	}
	viper.AutomaticEnv() // Automatically read environment variables
	return nil
}
