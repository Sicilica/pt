package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("push <tags>", commandPush, "Suspension & Swapping", "creates a new suspended task without starting it")
}

func commandPush(c types.CommandContext) error {
	tags := c.Args().Rest()
	if len(tags) == 0 {
		return errors.New("must specify one or more tags")
	}

	err := c.PT().PushSuspendedTags(tags)
	if err != nil {
		return err
	}

	fmt.Println("task pushed onto the stack")

	return nil
}
