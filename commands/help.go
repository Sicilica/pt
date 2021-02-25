package commands

import (
	"fmt"

	"github.com/sicilica/pt/types"
)

func init() {
	register("help", commandHelp, "Basic", "displays this help information")
}

func commandHelp(c types.CommandContext) error {
	for cat, cmds := range categoryCommands {
		fmt.Println()
		fmt.Println(cat)
		fmt.Println("------------")

		for _, cmd := range cmds {
			fmt.Println(commandExamples[cmd], "-", commandDescriptions[cmd])
		}
	}

	return nil
}
