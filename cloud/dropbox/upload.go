package dropbox

import (
	"bytes"
	"io"
	"os"

	"github.com/blend/go-sdk/crypto"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/pkg/errors"

	"github.com/sicilica/pt/cloud/dropbox/modifications"
)

func (d *dropboxCloudSyncProvider) Upload() error {
	info, err := os.Stat(d.localFile)
	if err != nil {
		return errors.Wrap(err, "failed to stat local file")
	}

	f, err := os.Open(d.localFile)
	if err != nil {
		return errors.Wrap(err, "failed to open local file")
	}
	defer f.Close()

	enc, err := crypto.NewStreamEncrypter(d.aesKey, f)
	if err != nil {
		return errors.Wrap(err, "failed to open encrypter")
	}

	if len(enc.IV) != encryptedBackupIVLen {
		return errors.New("unexpected iv length")
	}

	contents := io.MultiReader(
		bytes.NewReader(encryptedBackupMagicBytes),
		bytes.NewReader(enc.IV),
		enc,
	)

	dbx := dropbox.NewContext(d.cfg)
	_, err = modifications.FixedUpload(&dbx, &files.CommitInfo{
		Path:           d.remoteFile,
		Mode:           &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}},
		Mute:           true,
		ClientModified: info.ModTime(),
	}, contents)
	if err != nil {
		return errors.Wrap(err, "failed to upload file")
	}

	return nil
}
