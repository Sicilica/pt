package types

import (
	"github.com/pkg/errors"
)

// ErrNoOpenTask is returned by GetOpenTask if there is no open task.
var ErrNoOpenTask = errors.New("there is no open task")

// ErrNoParent is returned by GetTagParent if the tag has no parent.
var ErrNoParent = errors.New("tag has no parent")

// ErrNoSuspension is returned by PopSuspendedTags if there are no suspension
// records.
var ErrNoSuspension = errors.New("there is no suspended task")
