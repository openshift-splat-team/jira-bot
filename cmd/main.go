package main

import (
	"log"
	"os"

	"github.com/openshift-splat-team/splat-jira-bot/cmd/epic"
	"github.com/openshift-splat-team/splat-jira-bot/cmd/issue"
	"github.com/openshift-splat-team/splat-jira-bot/cmd/sprint"
	"github.com/openshift-splat-team/splat-jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

func initConfig() {
	util.BindEnvVars()
}

func main() {
	log.SetOutput(os.Stdout)
	initConfig()

	var rootCmd = &cobra.Command{Use: "jira-splat-bot"}

	epic.Initialize(rootCmd)
	sprint.Initialize(rootCmd)
	issue.Initialize(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
}
