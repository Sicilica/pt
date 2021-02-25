package util

import (
	"os"
	"path"

	"github.com/pkg/errors"
)

// GetLocalStorageDir calculates the local storage directory.
func GetLocalStorageDir() (string, error) {
	usrDir, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine local storage dir")
	}
	return path.Join(usrDir, "sicilica", "pt"), nil
}
