package util

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var bindEnvVars = []string{"personal_access_token", "project", "board"}

func CheckForMissingEnvVars() error {
	for _, envVar := range bindEnvVars {
		if len(viper.GetString(envVar)) == 0 {
			return fmt.Errorf("the environment variable: [%s] must be exported", strings.ToUpper(fmt.Sprintf("jira_%s", envVar)))
		}
	}
	return nil
}

func BindEnvVars() {
	viper.SetEnvPrefix("jira") // Set a prefix for environment variables
	for _, envVar := range bindEnvVars {
		viper.BindEnv(envVar)
	}
	viper.AutomaticEnv() // Automatically read environment variables
}
