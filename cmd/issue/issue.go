package issue

import (
	"github.com/spf13/cobra"
)

type issueCommandOptions struct {
	defaultSpikeStoryPoints int64
	points                  int64
	dryRunFlag              bool
	overrideFlag            bool
	state                   string
	priority                string
	comment 				string
	resolution              string
	summary                 string
	description             string
	issueType               string
	project                 string
}

var options = issueCommandOptions{
	defaultSpikeStoryPoints: -1,
	dryRunFlag:              true,
	overrideFlag:            false,
	priority:                "",
}

var cmdIssue = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Long:  `This command allows you to manage issues in your project management tool.`,
}

func Initialize(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cmdIssue)
}
