package util

import (
	"os"
)

// IsStatFileNotFound returns whether or not the given error represents a NOT_FOUND
// error from os.Stat.
func IsStatFileNotFound(err error) bool {
	if pe, ok := err.(*os.PathError); ok {
		if pe.Op == "CreateFile" {
			return true
		}
	}

	return false
}
