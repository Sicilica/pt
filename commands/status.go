package commands

import (
	"fmt"
	"time"

	"github.com/sicilica/pt/types"
)

func init() {
	register("status", commandStatus, "Basic", "displays information about the active and suspended tasks")
}

func commandStatus(c types.CommandContext) error {
	err := c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	t, err := c.PT().GetOpenTask()
	if err == types.ErrNoOpenTask {
		fmt.Println("idle")
	} else if err != nil {
		return err
	} else {
		tags, err := c.PT().GetTaskTags(t)
		if err != nil {
			return err
		}
		fmt.Printf("doing %s for %v since %s\n", formatTags(tags), formatDuration(time.Since(t.Start)), formatTime(t.Start))
	}

	suspensions, err := c.PT().ListSuspendedTasks()
	if err != nil {
		return err
	}

	if len(suspensions) > 0 {
		fmt.Println()
	}
	for i, sr := range suspensions {
		fmt.Printf("%d: %s @ %s\n", i, formatTags(sr.Tags), formatTime(sr.Time))
	}

	return nil
}
