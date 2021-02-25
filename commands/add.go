package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("add <tags>", commandAdd, "Tagging", "adds the given tags to the current task")
}

func commandAdd(c types.CommandContext) error {
	tags := c.Args().Rest()
	if len(tags) == 0 {
		return errors.New("must specify one or more tags")
	}

	t, err := c.PT().GetOpenTask()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := c.PT().AddTagToTask(t, tag); err != nil {
			return err
		}
	}

	fmt.Println("tags added")

	return nil
}
