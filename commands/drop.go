package commands

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("drop [index]", commandDrop, "Suspension & Swapping", "drops a suspended task from the stack")
}

func commandDrop(c types.CommandContext) error {
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

	fmt.Println("dropped suspended task with tags", formatTags(tags))

	return nil
}
