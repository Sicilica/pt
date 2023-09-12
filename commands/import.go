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
	register("import <file>", commandImport, "Storage", "import database from a file")
}

func commandImport(c types.CommandContext) error {
	err := c.CloseStorage()
	if err != nil {
		return errors.Wrap(err, "failed to close storage")
	}

	appDir, err := util.GetLocalStorageDir()
	if err != nil {
		return err
	}

	importFile, err := c.Args().Pop()
	if err != nil {
		return err
	}
	err = c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	in, err := os.Open(importFile)
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

	fmt.Println("data imported from file", importFile)

	return nil
}
