package commands

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("rollback [start time] <tags>", commandRollback, "Temporal", "starts a new task at a set point in the past, or when the last task ended")
}

func commandRollback(c types.CommandContext) error {
	// Parse start time arg
	useLastStopTime := false
	startTime, err := c.Args().PopTime(time.Now(), false)
	if err != nil {
		useLastStopTime = true
	} else if !startTime.Before(time.Now()) {
		return errors.New("start time must be in the past")
	}

	tags := c.Args().Rest()
	if len(tags) == 0 {
		return errors.New("must specify one or more tags")
	}

	t, err := c.PT().GetOpenTask()
	if err != types.ErrNoOpenTask {
		if err != nil {
			return errors.Wrap(err, "failed to close existing task")
		}

		if useLastStopTime {
			return errors.New("a task is already open, so a start time is required")
		}

		if err := c.PT().SetStopTime(t, startTime); err != nil {
			return errors.Wrap(err, "failed to close existing task")
		}
	}

	recent, err := c.PT().GetLastTask()
	if err != nil {
		return err
	}
	if useLastStopTime {
		startTime = recent.Stop
	} else if startTime.Before(recent.Stop) {
		return errors.New("start time cannot overlap with previous task")
	}

	t, err = c.PT().CreateTask(startTime)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := c.PT().AddTagToTask(t, tag); err != nil {
			return err
		}
	}

	fmt.Println("task started at", formatTime(t.Start))

	return nil
}
