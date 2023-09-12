package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sicilica/pt/types"
)

func init() {
	register("i", commandI, "Basic", "run pt in interactive mode")
}

func commandI(c types.CommandContext) error {
	if err := c.Args().MustBeEmpty(); err != nil {
		return err
	}

	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !s.Scan() {
			break
		}

		// TODO This isn't a perfect parser by any means, but I'm too lazy to import
		// anything for this right now
		args := strings.Split(s.Text(), " ")
		for i, a := range args {
			args[i] = strings.TrimSpace(a)
		}

		if err := c.ExecCommand(args); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	fmt.Println()

	return nil
}
