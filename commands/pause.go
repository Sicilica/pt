package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("pause", commandPause, "Suspension & Swapping", "suspends the current task")
}

func commandPause(c types.CommandContext) error {
	err := c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	t, err := mustCloseExistingTask(c)
	if err != nil {
		return err
	}

	tags, err := c.PT().GetTaskTags(t)
	if err != nil {
		return errors.Wrap(err, "failed to get task tags")
	}

	err = c.PT().PushSuspendedTags(tags)
	if err != nil {
		return errors.Wrap(err, "failed to create suspension record")
	}

	fmt.Println("task suspended at", formatTime(t.Stop))

	return nil
}
