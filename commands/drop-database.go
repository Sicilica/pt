package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
	"github.com/sicilica/pt/util"
)

func init() {
	register("drop-database", commandDropDatabase, "Storage", "deletes the entire database (are you sure about this?)")
}

func commandDropDatabase(c types.CommandContext) error {
	err := c.CloseStorage()
	if err != nil {
		return errors.Wrap(err, "failed to close storage")
	}

	appDir, err := util.GetLocalStorageDir()
	if err != nil {
		return err
	}

	if err := os.Remove(path.Join(appDir, "pt.db")); err != nil {
		return err
	}

	fmt.Println("database deleted (backups may still exist)")

	return nil
}
