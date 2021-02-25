package commands

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("out", commandOut, "Basic", "ends the current task")
}

func commandOut(c types.CommandContext) error {
	err := c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	t, err := mustCloseExistingTask(c)
	if err != nil {
		return err
	}

	fmt.Println("task stopped at", formatTime(t.Stop))

	return nil
}

func mustCloseExistingTask(c types.CommandContext) (*types.Task, error) {
	t, err := c.PT().GetOpenTask()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close task")
	}

	err = c.PT().SetStopTime(t, time.Now())
	if err != nil {
		return nil, errors.Wrap(err, "failed to close task")
	}

	return t, nil
}
