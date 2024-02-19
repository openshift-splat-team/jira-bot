package sprint

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/splat-jira-bot/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func moveToSprint(client *jira.Client, sprintNumber string, issues []string) error {
	log.Printf("moving %d issues to sprint %s\n", len(issues), sprintNumber)
	boards, _, err := client.Board.GetAllBoards(&jira.BoardListOptions{
		ProjectKeyOrID: viper.GetString("project"),
	})
	if err != nil {
		return fmt.Errorf("unable to get boards: %v", err)
	}

	targetBoard := viper.GetString("board")
	log.Printf("finding board %s\n", targetBoard)

	// "SPLAT - Scrum Board"
	var board *jira.Board
	for _, _board := range boards.Values {
		if _board.Name == targetBoard {
			board = &_board
			break
		}
	}
	if board == nil {
		return fmt.Errorf("unable to find board %s", targetBoard)
	}

	log.Printf("found board %s; id: %d\n", targetBoard, board.ID)

	log.Printf("finding sprint %s\n", sprintNumber)
	sprints, _, err := client.Board.GetAllSprints(strconv.Itoa(board.ID))

	if err != nil {
		return fmt.Errorf("unable to get sprint: %v", err)
	}

	var sprint *jira.Sprint
	for _, _sprint := range sprints {
		if _sprint.Name == sprintNumber {
			sprint = &_sprint
			break
		}
	}

	if sprint == nil {
		return fmt.Errorf("unable to find sprint %s", sprintNumber)
	}

	log.Printf("found sprint %s\n", sprintNumber)

	if !dryRunFlag {
		log.Printf("moving issues in to sprint\n")
		_, err = client.Sprint.MoveIssuesToSprint(sprint.ID, issues)
		if err != nil {
			return fmt.Errorf("unable to move issues to sprint %s: %v", sprintNumber, err)
		}
	} else {
		log.Printf("issues: %s would have been moved to sprint %s. run again and provide --dry-run=false to apply.", strings.Join(issues, ","), sprintNumber)
	}

	return err
}

var cmdMoveIssue = &cobra.Command{
	Use:   "move-issue [sprint-number] [issue-number]...",
	Short: "Move an issue to a sprint",
	Long:  `This command allows you to move one or more issues to a sprint in your project management tool.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := util.CheckForMissingEnvVars()
		if err != nil {
			util.RuntimeError(err)
		}
		sprintNumber := args[0]
		issueNumbers := args[1:]
		client, err := util.GetJiraClient()
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to get jira client: %v", err))
		}
		err = moveToSprint(client, sprintNumber, issueNumbers)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to move issue: %v", err))
		}
	},
}
