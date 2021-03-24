package commands

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("out [end time]", commandOut, "Basic", "ends the current task")
}

func commandOut(c types.CommandContext) error {
	// Parse end time arg
	endTime, err := c.Args().PopTime(time.Now(), false)
	if err != nil {
		endTime = time.Now()
	} else if !endTime.Before(time.Now()) {
		return errors.New("end time must be in the past")
	}

	err = c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	t, err := mustCloseExistingTask(c, endTime)
	if err != nil {
		return err
	}

	fmt.Println("task stopped at", formatTime(t.Stop))

	return nil
}

func mustCloseExistingTask(c types.CommandContext, endTime time.Time) (*types.Task, error) {
	t, err := c.PT().GetOpenTask()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close task")
	}

	err = c.PT().SetStopTime(t, endTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to close task")
	}

	return t, nil
}
