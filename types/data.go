package types

import (
	"time"
)

// Query fields can be used to further limit task results.
type Query struct {
	Tag string
}

// A Task is clock entry that has a start and stop time and is attached to tags.
type Task struct {
	ID    string
	Start time.Time
	Stop  time.Time
}

// A SuspensionRecord is a stack entry that tracks the state of a suspended task.
type SuspensionRecord struct {
	Tags []string
	Time time.Time
}
