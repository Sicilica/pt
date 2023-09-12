package app

import (
	"github.com/pkg/errors"

	"github.com/sicilica/pt/commands"
	"github.com/sicilica/pt/types"
	"github.com/sicilica/pt/util"
)

type Runtime struct {
	NewStorageProvider   func() (types.StorageProvider, error)
	NewCloudSyncProvider func() (types.CloudSyncProvider, error)

	execDepth int
	storageProvider types.StorageProvider
}

func (r *Runtime) Close() error {
	if r.storageProvider != nil {
		defer r.storageProvider.Close()
	}

	return nil
}

func (r *Runtime) ExecCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("no command specified")
	}

	c := &commandContext{
		args:    util.NewArgsQueue(args),
		runtime: r,
	}
	defer c.cleanup()

	cmdName, err := c.Args().Pop()
	if err != nil {
		return err
	}

	cmd, ok := commands.Get(cmdName)
	if !ok {
		return errors.Errorf("unrecognized command \"%s\"", cmdName)
	}

	r.execDepth++
	err = cmd(c)
	r.execDepth--
	if err != nil {
		return err
	}

	if c.storageSession != nil {
		err = c.storageSession.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}
