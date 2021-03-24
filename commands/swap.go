package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("swap [index|tags]", commandSwap, "Suspension & Swapping", "suspends the current task and replaces it with either a new task or a previously suspended task")
}

func commandSwap(c types.CommandContext) error {
	// Figure out the tags we're using, and pop the suspended task if necessary.
	var newTags []string
	maybeIdxS := c.Args().Peek()
	_, err := strconv.Atoi(maybeIdxS)
	if maybeIdxS == "" || err == nil {
		idx := 0
		if err == nil {
			c.Args().MustPop(maybeIdxS)
			if err != nil {
				return err
			}
			idx, err = strconv.Atoi(maybeIdxS)
			if err != nil {
				return err
			}
			if idx < 0 {
				return errors.New("invalid index")
			}
		}

		err = c.Args().MustBeEmpty()
		if err != nil {
			return err
		}

		newTags, err = c.PT().PopSuspendedTags(idx)
		if err == types.ErrNoSuspension {
			if idx > 0 {
				return errors.New("index out of range")
			}
			return errors.New("no tags specified, and no suspended task found")
		} else if err != nil {
			return err
		}
	} else {
		newTags = c.Args().Rest()
		if len(newTags) == 0 {
			return errors.New("PANIC! assumed tag-based, but no tags specified")
		}
	}

	// Now we essentially do a normal "resume" call...
	prevT, err := mustCloseExistingTask(c, time.Now())
	if err != nil {
		return err
	}

	prevTags, err := c.PT().GetTaskTags(prevT)
	if err != nil {
		return errors.Wrap(err, "failed to get task tags")
	}

	err = c.PT().PushSuspendedTags(prevTags)
	if err != nil {
		return errors.Wrap(err, "failed to create suspension record")
	}

	// And finally, a normal "in" call.
	newT, err := c.PT().CreateTask(time.Now())
	if err != nil {
		return err
	}

	for _, tag := range newTags {
		if err := c.PT().AddTagToTask(newT, tag); err != nil {
			return err
		}
	}

	fmt.Printf("task swapped from %s to %s at %s\n", formatTags(prevTags), formatTags(newTags), formatTime(newT.Start))

	return nil
}
