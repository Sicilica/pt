package commands

import (
	"fmt"

	"github.com/sicilica/pt/types"
)

func init() {
	register("forget <tag>", commandForget, "Tagging", "removes the parent from a given tag")
}

func commandForget(c types.CommandContext) error {
	tag, err := c.Args().Pop()
	if err != nil {
		return err
	}
	c.Args().MustBeEmpty()

	err = c.PT().ClearTagParent(tag)
	if err != nil {
		return err
	}

	fmt.Printf("\"%s\" has no parent\n", tag)

	return nil
}
