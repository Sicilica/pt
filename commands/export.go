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
	register("export <file>", commandExport, "Storage", "export database to a file")
}

func commandExport(c types.CommandContext) error {
	err := c.CloseStorage()
	if err != nil {
		return errors.Wrap(err, "failed to close storage")
	}

	appDir, err := util.GetLocalStorageDir()
	if err != nil {
		return err
	}

	exportFile, err := c.Args().Pop()
	if err != nil {
		return err
	}
	err = c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	in, err := os.Open(path.Join(appDir, "pt.db"))
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(exportFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	fmt.Println("data exported to file", exportFile)

	return nil
}
