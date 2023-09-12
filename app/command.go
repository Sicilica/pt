package app

import (
	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
	"github.com/sicilica/pt/util"
)

type commandContext struct {
	args           *util.ArgsQueue
	storageSession types.StorageSession
	runtime        *Runtime
}

var _ types.CommandContext = (*commandContext)(nil)

func (c *commandContext) cleanup() error {
	if c.storageSession != nil {
		defer c.storageSession.Abort()
	}

	return nil
}

func (c *commandContext) Args() *util.ArgsQueue {
	return c.args
}

func (c *commandContext) CloseStorage() error {
	if c.storageSession != nil {
		return errors.New("command has a storage session")
	}

	if c.runtime.storageProvider != nil {
		err := c.runtime.storageProvider.Close()
		if err != nil {
			return err
		}
		c.runtime.storageProvider = nil
	}

	return nil
}

func (c *commandContext) DisallowNesting() error {
	if c.runtime.execDepth > 1 {
		return errors.New("command cannot be executed in interactive mode")
	}
	return nil
}

func (c *commandContext) ExecCommand(args []string) error {
	if c.storageSession != nil {
		return errors.New("command has a storage session and cannot spawn children")
	}

	return c.runtime.ExecCommand(args)
}

func (c *commandContext) PT() types.StorageInterface {
	var err error
	if c.storageSession == nil {
		if c.runtime.storageProvider == nil {
			c.runtime.storageProvider, err = c.runtime.NewStorageProvider()
			if err != nil {
				panic(err)
			}
		}

		c.storageSession, err = c.runtime.storageProvider.NewSession()
		if err != nil {
			panic(err)
		}
	}

	return c.storageSession
}

func (c *commandContext) WithCloudSync(fn func(cloud types.CloudSyncInterface) error) error {
	if err := c.CloseStorage(); err != nil {
		return errors.Wrap(err, "failed to close storage")
	}

	cloud, err := c.runtime.NewCloudSyncProvider()
	if err != nil {
		return errors.Wrap(err, "failed to init cloud sync provider")
	}
	defer cloud.Close()

	return fn(cloud)
}
