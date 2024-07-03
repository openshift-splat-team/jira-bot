package sprint

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cmdSyncFutureSprints.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")

	cmdSprint.AddCommand(cmdSyncFutureSprints)
}

func createFutureSprints(ctx context.Context, client *jira.Client, futureSprints int64) error {
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

	log.Println("finding current, active sprint")
	sprints, _, err := client.Board.GetAllSprints(strconv.Itoa(board.ID))

	if err != nil {
		return fmt.Errorf("unable to get sprint: %v", err)
	}

	var sprint *jira.Sprint
	for _, _sprint := range sprints {
		if _sprint.State == "active" {
			sprint = &_sprint
			break
		}
	}

	if sprint == nil {
		return fmt.Errorf("unable to find active sprint")
	}

	log.Printf("found sprint %s\n", sprint.Name)

	nameParts := strings.Split(sprint.Name, " ")
	slices.Reverse(nameParts)

	activeSprintNumber, err := strconv.Atoi(nameParts[0])
	if err != nil {
		return fmt.Errorf("unable to convert active sprint number to int: %v", err)
	}

	for i := 1; i <= int(futureSprints); i++ {
		sprintName := fmt.Sprintf("OpenShift SPLAT - Sprint %d", activeSprintNumber+i)
		skip := false
		for _, knownScript := range sprints {
			if knownScript.Name == sprintName {
				log.Printf("sprint %s already exists\n", sprintName)
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		startDate := sprint.StartDate.AddDate(0, 0, 21*i)
		endDate := sprint.EndDate.AddDate(0, 0, 21*i)
		newSprint := &jira.Sprint{
			Name:          sprintName,
			OriginBoardID: board.ID,
			StartDate:     &startDate,
			EndDate:       &endDate,
		}
		if !dryRunFlag {

			_, resp, err := client.Sprint.CreateSprint(ctx, newSprint)
			if err != nil {
				responseBody, _ := util.GetResponseBody(resp)
				return fmt.Errorf("unable to create sprint: %v. response body: %s", err, responseBody)
			}
			log.Printf("created sprint %s\n", newSprint.Name)
		} else {
			log.Printf("sprint %s would have been created as: %v. run again and provide --dry-run=false to apply.", sprintName, newSprint)
		}
	}

	return err
}

var cmdSyncFutureSprints = &cobra.Command{
	Use:   "sync-future-sprints [number of sprints from active]",
	Short: "Precreates future sprints to aid in sprint planning",
	Long:  `This command allows you to precreate future sprints to aid in sprint planning`,
	Args:  cobra.ExactArgs(1), // Requires exactly one argument: sprint-number and issue-number
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()
		err := util.CheckForMissingEnvVars()
		if err != nil {
			util.RuntimeError(err)
		}
		numberOfSprints := args[0]
		client, err := util.GetJiraClient()
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to get jira client: %v", err))
		}
		sprints, err := strconv.Atoi(numberOfSprints)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to convert number of sprints to int: %v", err))
		}
		err = createFutureSprints(ctx, client, int64(sprints))
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to create future sprints: %v", err))
		}
	},
}
