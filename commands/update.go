package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/sicilica/pt/types"
)

func init() {
	register("update", commandUpdate, "Basic", "update pt itself to the latest version")
}

func commandUpdate(c types.CommandContext) error {
	err := c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	if err := os.Chdir(path.Dir(path.Dir(exe))); err != nil {
		return err
	}

	cmd := exec.Command("git", "pull")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return err
	}
	if buf.String() == "Already up to date.\n" {
		fmt.Println("already up to date")
		return nil
	}

	if err := exec.Command("make").Run(); err != nil {
		return err
	}

	fmt.Println("updated successfully")

	return nil
}
