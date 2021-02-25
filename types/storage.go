package types

import (
	"time"
)

// StorageInterface is the interface of endpoints that storage providers must expose.
type StorageInterface interface {
	AddTagToTask(t *Task, tag string) error
	ClearTagParent(tag string) error
	CreateTask(start time.Time) (*Task, error)
	GetLastTask() (*Task, error)
	GetOpenTask() (*Task, error)
	GetTagParent(tag string) (string, error)
	GetTaskTags(t *Task) ([]string, error)
	ListSuspendedTasks() ([]*SuspensionRecord, error)
	PopSuspendedTags(offset int) ([]string, error)
	PushSuspendedTags(tags []string) error
	QueryTasks(start time.Time, end time.Time, q *Query) ([]*Task, error)
	RemoveTagFromTask(t *Task, tag string) error
	SetStartTime(t *Task, start time.Time) error
	SetStopTime(t *Task, stop time.Time) error
	SetTagParent(child, parent string) error
	TagHasAncestor(child, parent string) (bool, error)
}

// CloudSyncInterface is the interface of endpoints that cloud sync providers must
// expose.
type CloudSyncInterface interface {
	Download() error
	GetLocalVersion() (BackupVersionInfo, error)
	GetRemoteVersion() (BackupVersionInfo, error)
	Upload() error
}

// BackupVersionInfo describes the metadata about a specific backup.
type BackupVersionInfo struct {
	Exists bool
	Modified time.Time
}

// A StorageProvider is an active connection to a datastore. Actual access to the
// storage is protected by sessions to allow easy rolling back.
type StorageProvider interface {
	Close() error
	NewSession() (StorageSession, error)
}

// A StorageSession is an active connection to a datastore, scoping to the execution
// of a single command.
type StorageSession interface {
	StorageInterface
	Abort() error
	Commit() error
}

// A CloudSyncProvider is an active cloud sync connection.
type CloudSyncProvider interface {
	CloudSyncInterface
	Close() error
}
