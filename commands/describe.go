package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sicilica/pt/types"
)

func init() {
	register("describe <child> [parent]", commandDescribe, "Tagging", "configures a given tag to be the child of another tag, or displays its current parent")
}

func commandDescribe(c types.CommandContext) error {
	child, err := c.Args().Pop()
	if err != nil {
		return err
	}

	if parent, err := c.Args().Pop(); err == nil {
		hasAncestor, err := c.PT().TagHasAncestor(parent, child)
		if err != nil {
			return err
		}
		if hasAncestor {
			return errors.New("circular hierarchy detected")
		}

		err = c.PT().SetTagParent(child, parent)
		if err != nil {
			return err
		}

		fmt.Printf("\"%s\" is a child of \"%s\"\n", child, parent)
	} else {
		parent, err := c.PT().GetTagParent(child)
		if err == types.ErrNoParent {
			fmt.Printf("\"%s\" has no parent\n", child)
		} else if err != nil {
			return err
		} else {
			fmt.Printf("\"%s\" is a child of \"%s\"\n", child, parent)
		}
	}

	return nil
}
