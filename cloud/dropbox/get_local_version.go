package dropbox

import (
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
	"github.com/sicilica/pt/util"
)

func (d *dropboxCloudSyncProvider) GetLocalVersion() (types.BackupVersionInfo, error) {
	info, err := os.Stat(d.localFile)
	if err != nil {
		if util.IsStatFileNotFound(err) {
			return types.BackupVersionInfo{
				Exists: false,
			}, nil
		}
		return types.BackupVersionInfo{}, errors.Wrap(err, "failed to stat local file")
	}

	return types.BackupVersionInfo{
		Exists:   true,
		Modified: info.ModTime().UTC().Truncate(time.Second),
	}, nil
}
