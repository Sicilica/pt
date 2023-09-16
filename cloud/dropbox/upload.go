package dropbox

import (
	"bytes"
	"io"
	"os"

	"github.com/blend/go-sdk/crypto"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
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

	cfg, err := d.Config()
	if err != nil {
		return err
	}

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

	modTime := info.ModTime()

	dbx := dropbox.NewContext(cfg)
	err = modifications.FixedUpload(&dbx, &files.CommitInfo{
		Path:           d.remoteFile,
		Mode:           &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}},
		Mute:           true,
		ClientModified: &modTime,
	}, contents)
	if err != nil {
		return errors.Wrap(err, "failed to upload file")
	}

	return nil
}
