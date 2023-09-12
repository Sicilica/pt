package util

import (
	"os"
	"syscall"
)

// IsStatFileNotFound returns whether or not the given error represents a NOT_FOUND
// error from os.Stat.
func IsStatFileNotFound(err error) bool {
	if pe, ok := err.(*os.PathError); ok {
		if errno, ok := pe.Err.(syscall.Errno); ok {
			return errno == syscall.ENOENT
		}
	}

	return false
}
