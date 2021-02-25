package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("resume [index]", commandResume, "Suspension & Swapping", "resumes a suspended task")
}

func commandResume(c types.CommandContext) error {
	idx := 0
	idxS, err := c.Args().Pop()
	if err == nil {
		idx, err = strconv.Atoi(idxS)
		if err != nil {
			return errors.Wrap(err, "failed to parse index")
		}
		if idx < 0 {
			return errors.New("invalid index")
		}
	}

	tags, err := c.PT().PopSuspendedTags(idx)
	if err == types.ErrNoSuspension && idx > 0 {
		return errors.New("index out of range")
	} else if err != nil {
		return err
	}

	err = autoCloseOpenTask(c)
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

	fmt.Println("task resumed at", formatTime(t.Start))

	return nil
}
