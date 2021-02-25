package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("remove <tags>", commandRemove, "Tagging", "removes the given tags from the current task")
}

func commandRemove(c types.CommandContext) error {
	tags := c.Args().Rest()
	if len(tags) == 0 {
		return errors.New("must specify one or more tags")
	}

	t, err := c.PT().GetOpenTask()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := c.PT().RemoveTagFromTask(t, tag); err != nil {
			return err
		}
	}

	tags, err = c.PT().GetTaskTags(t)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return errors.New("cannot remove all tags from a task")
	}

	fmt.Println("tags removed")

	return nil
}
