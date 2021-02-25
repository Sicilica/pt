package dropbox

import (
	"encoding/base64"
	"encoding/hex"
	"path"

	"github.com/blend/go-sdk/crypto"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
	"github.com/sicilica/pt/util"
)

const encryptedBackupIVLen = 16

var encryptedBackupMagicBytes []byte

func init() {
	var err error
	encryptedBackupMagicBytes, err = hex.DecodeString("90beef77")
	if err != nil {
		panic(errors.Wrap(err, "failed to generate magic bytes"))
	}
}

type dropboxCloudSyncProvider struct {
	cfg        dropbox.Config
	aesKey     []byte
	localFile  string
	remoteFile string
}

// sqlite3Provider implements StorageProvider.
var _ types.CloudSyncProvider = (*dropboxCloudSyncProvider)(nil)

// New returns a dropbox-powered sync storage interface.
func New() (types.CloudSyncProvider, error) {
	appDir, err := util.GetLocalStorageDir()
	if err != nil {
		return nil, err
	}
	cfgFile := path.Join(appDir, "dropbox_config.json")

	cfg := struct {
		AccessToken   string `json:"access_token"`
		EncryptionKey string `json:"encryption_key"`
	}{}
	err = util.ReadJSONFile(cfgFile, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read dropbox config")
	}

	if cfg.AccessToken == "" {
		return nil, errors.New("missing required access_token in dropbox config")
	}

	var encryptionKey []byte
	if cfg.EncryptionKey != "" {
		encryptionKey, err = base64.StdEncoding.DecodeString(cfg.EncryptionKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse encryption key")
		}
	} else {
		encryptionKey, err = crypto.CreateKey(32)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate encryption key")
		}
		cfg.EncryptionKey = base64.StdEncoding.EncodeToString(encryptionKey)

		err = util.WriteJSONFile(cfgFile, &cfg, 0600)
		if err != nil {
			return nil, errors.Wrap(err, "failed to write updated dropbox config")
		}
	}

	if len(encryptionKey) != 32 {
		return nil, errors.New("encryption key has invalid length")
	}

	return &dropboxCloudSyncProvider{
		cfg: dropbox.Config{
			Token: cfg.AccessToken,
		},
		aesKey:     encryptionKey,
		localFile:  path.Join(appDir, "pt.db"),
		remoteFile: "/pt.db",
	}, nil
}

func (d *dropboxCloudSyncProvider) Close() error {
	return nil
}
