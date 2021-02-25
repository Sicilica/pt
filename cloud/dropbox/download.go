package dropbox

import (
	"io"
	"os"

	"github.com/blend/go-sdk/crypto"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/pkg/errors"
)

func (d *dropboxCloudSyncProvider) Download() error {
	_, contents, err := files.New(d.cfg).Download(files.NewDownloadArg(d.remoteFile))
	if err != nil {
		return errors.Wrap(err, "failed to download file")
	}
	defer contents.Close()

	magicBytes := make([]byte, len(encryptedBackupMagicBytes))
	count, err := contents.Read(magicBytes)
	if err != nil {
		return errors.Wrap(err, "error while reading magic bytes")
	}
	if count != len(magicBytes) {
		return errors.New("read wrong count while reading magic bytes")
	}
	for i := 0; i < len(encryptedBackupMagicBytes); i++ {
		if magicBytes[i] != encryptedBackupMagicBytes[i] {
			return errors.New("magic bytes didn't match")
		}
	}

	iv := make([]byte, encryptedBackupIVLen)
	count, err = contents.Read(iv)
	if err != nil {
		return errors.Wrap(err, "error while reading iv")
	}
	if count != len(iv) {
		return errors.New("read wrong count while reading iv")
	}

	dec, err := crypto.NewStreamDecrypter(d.aesKey, crypto.StreamMeta{
		IV: iv,
	}, contents)
	if err != nil {
		return errors.Wrap(err, "failed to open decrypter")
	}

	f, err := os.Create(d.localFile)
	if err != nil {
		return errors.Wrap(err, "failed to create local file")
	}
	defer f.Close()

	_, err = io.Copy(f, dec)
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	return nil
}
