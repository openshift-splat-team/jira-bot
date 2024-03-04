package main

import (
	"fmt"
	"log"
	"os"

	"github.com/openshift-splat-team/jira-bot/cmd/epic"
	"github.com/openshift-splat-team/jira-bot/cmd/issue"
	"github.com/openshift-splat-team/jira-bot/cmd/sprint"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

func initConfig() error {
	return util.BindEnvVars()
}

func main() {
	log.SetOutput(os.Stdout)
	err := initConfig()
	if err != nil {
		util.RuntimeError(fmt.Errorf("unable to initialize: %v", err))
	}

	var rootCmd = &cobra.Command{Use: "jira-splat-bot"}

	epic.Initialize(rootCmd)
	sprint.Initialize(rootCmd)
	issue.Initialize(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		util.RuntimeError(fmt.Errorf("error while running command: %v", err))
	}
}
