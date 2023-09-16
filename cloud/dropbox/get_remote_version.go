package dropbox

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func (d *dropboxCloudSyncProvider) GetRemoteVersion() (types.BackupVersionInfo, error) {
	cfg, err := d.Config()
	if err != nil {
		return types.BackupVersionInfo{}, err
	}
	res, err := files.New(cfg).GetMetadata(files.NewGetMetadataArg(d.remoteFile))
	if err != nil {
		if gma, ok := err.(files.GetMetadataAPIError); ok {
			if gma.EndpointError.Path.Tag == "not_found" {
				return types.BackupVersionInfo{
					Exists: false,
				}, nil
			}
		}
		return types.BackupVersionInfo{}, errors.Wrap(err, "failed to get remote metadata")
	}

	fm, ok := res.(*files.FileMetadata)
	if !ok {
		return types.BackupVersionInfo{}, errors.New("failed to convert received metadata")
	}

	return types.BackupVersionInfo{
		Exists:   true,
		Modified: fm.ClientModified,
	}, nil
}
