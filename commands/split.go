package commands

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("split [split time] [add|remove <tags>]", commandSplit, "Temporal", "ends the current task and starts a copy of it")
}

func commandSplit(c types.CommandContext) error {
	// Parse split time arg
	splitTime := time.Now()
	maybeTime := c.Args().Peek()
	switch maybeTime {
	case "":
	case "add":
	case "remove":
		// noop
	default:
		var err error
		splitTime, err = c.Args().PopTime(time.Now(), false)
		if err != nil {
			return err
		}
		if !splitTime.Before(time.Now()) {
			return errors.New("split time must be in the past")
		}
	}

	// Parse add/remove tags
	var tagsToAdd []string
	var tagsToRemove []string
	switch c.Args().Peek() {
	case "add":
		c.Args().MustPop("add")
		tagsToAdd = c.Args().Rest()
	case "remove":
		c.Args().MustPop("remove")
		tagsToRemove = c.Args().Rest()
	default:
		if err := c.Args().MustBeEmpty(); err != nil {
			return err
		}
	}

	openT, err := c.PT().GetOpenTask()
	if err != nil {
		return err
	}

	if !splitTime.After(openT.Start) {
		return errors.New("split time must be after the start of the current task")
	}

	prevTags, err := c.PT().GetTaskTags(openT)
	if err != nil {
		return errors.Wrap(err, "failed to get task tags")
	}

	// Verification step only: make sure all tags were set
	for _, tag := range tagsToRemove {
		found := false
		for _, prevTag := range prevTags {
			if tag == prevTag {
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("tag \"%s\" is not currently active", tag)
		}
	}

	splitT, err := c.PT().CreateTask(openT.Start)
	if err != nil {
		return errors.Wrap(err, "failed to create new task")
	}

	err = c.PT().SetStopTime(splitT, splitTime)
	if err != nil {
		return errors.Wrap(err, "failed to close new task")
	}

	for _, tag := range prevTags {
		if err := c.PT().AddTagToTask(splitT, tag); err != nil {
			return errors.Wrap(err, "failed to set new task's tags")
		}
	}

	err = c.PT().SetStartTime(openT, splitTime)
	if err != nil {
		return errors.Wrap(err, "failed to update open task start")
	}

	for _, tag := range tagsToAdd {
		if err := c.PT().AddTagToTask(openT, tag); err != nil {
			return err
		}
	}
	for _, tag := range tagsToRemove {
		if err := c.PT().RemoveTagFromTask(openT, tag); err != nil {
			return err
		}
	}

	tags, err := c.PT().GetTaskTags(openT)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return errors.New("cannot remove all tags from a task")
	}

	fmt.Println("task split at", formatTime(splitTime))

	return nil
}
