package commands

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
	"github.com/sicilica/pt/util"
)

func init() {
	register("restore [name]", commandRestore, "Storage", "restores a previously created backup")
}

func commandRestore(c types.CommandContext) error {
	err := c.CloseStorage()
	if err != nil {
		return errors.Wrap(err, "failed to close storage")
	}

	appDir, err := util.GetLocalStorageDir()
	if err != nil {
		return err
	}

	backupFile := "pt-backup.db"
	name, err := c.Args().Pop()
	if err == nil {
		backupFile = fmt.Sprintf("pt-backup-%s.db", name)
	}
	err = c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	in, err := os.Open(path.Join(appDir, backupFile))
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(path.Join(appDir, "pt.db"))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	fmt.Println("restored backup from disk")

	return nil
}
