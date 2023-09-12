package types

import (
	"github.com/sicilica/pt/util"
)

// CommandContext is the main object passed to running commands and wraps their
// state.
type CommandContext interface {
	Args() *util.ArgsQueue
	CloseStorage() error
	DisallowNesting() error
	ExecCommand(args []string) error
	PT() StorageInterface
	WithCloudSync(fn func(cloud CloudSyncInterface) error) error
}
