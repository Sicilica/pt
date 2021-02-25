package commands

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("in <tags>", commandIn, "Basic", "starts a new task, ending the current task if one exists")
}

func commandIn(c types.CommandContext) error {
	tags := c.Args().Rest()
	if len(tags) == 0 {
		return errors.New("must specify one or more tags")
	}

	err := autoCloseOpenTask(c)
	if err != nil {
		return err
	}

	t, err := c.PT().CreateTask(time.Now())
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

func autoCloseOpenTask(c types.CommandContext) error {
	t, err := c.PT().GetOpenTask()
	if err == types.ErrNoOpenTask {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to close existing task")
	}

	if err := c.PT().SetStopTime(t, time.Now()); err != nil {
		return errors.Wrap(err, "failed to close existing task")
	}
	fmt.Println("closing existing task")
	return nil
}
