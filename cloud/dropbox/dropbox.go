package dropbox

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/blend/go-sdk/crypto"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/pkg/errors"
	"github.com/sicilica/pt/types"
	"github.com/sicilica/pt/util"
)

const defaultClientID = "xhkhup08ziqudtq"

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
	accessToken          string
	accessTokenExpiresAt time.Time
	aesKey               []byte
	config               dropboxConfig
	localFile            string
	remoteFile           string
}

type dropboxConfig struct {
	ClientID      string `json:"client_id,omitempty"`
	EncryptionKey string `json:"encryption_key"`
	RefreshToken  string `json:"refresh_token"`
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

	cfg := dropboxConfig{}
	_, err = os.Stat(cfgFile)
	if err == nil {
		if err := util.ReadJSONFile(cfgFile, &cfg); err != nil {
			return nil, err
		}
	} else if !util.IsStatFileNotFound(err) {
		return nil, err
	}

	var encryptionKey []byte
	if cfg.EncryptionKey != "" {
		encryptionKey, err = base64.StdEncoding.DecodeString(cfg.EncryptionKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse encryption key")
		}
	}

	if len(encryptionKey) != 32 {
		return nil, errors.New("encryption key has invalid length")
	}

	return &dropboxCloudSyncProvider{
		aesKey:     encryptionKey,
		config:     cfg,
		localFile:  path.Join(appDir, "pt.db"),
		remoteFile: "/pt.db",
	}, nil
}

func (d *dropboxCloudSyncProvider) Close() error {
	return nil
}

func (d *dropboxCloudSyncProvider) Config() (dropbox.Config, error) {
	tok, err := d.AccessToken()
	return dropbox.Config{
		Token: tok,
	}, err
}

func (d *dropboxCloudSyncProvider) AccessToken() (string, error) {
	configUpdated := false
	defer func() {
		if configUpdated {
			appDir, err := util.GetLocalStorageDir()
			if err != nil {
				panic(err)
			}
			cfgFile := path.Join(appDir, "dropbox_config.json")

			err = util.WriteJSONFile(cfgFile, &d.config, 0600)
			if err != nil {
				panic(errors.Wrap(err, "failed to write updated dropbox config"))
			}
		}
	}()

	var err error

	if d.config.EncryptionKey == "" {
		d.aesKey, err = crypto.CreateKey(32)
		if err != nil {
			return "", errors.Wrap(err, "failed to generate encryption key")
		}
		d.config.EncryptionKey = base64.StdEncoding.EncodeToString(d.aesKey)
		configUpdated = true
	}

	if d.accessTokenExpiresAt.After(time.Now().Add(15 * time.Second)) {
		return d.accessToken, nil
	}

	clientID := d.config.ClientID
	if clientID == "" {
		clientID = defaultClientID
	}

	var tokenReq url.Values
	if d.config.RefreshToken == "" {
		// Generate a code verifier
		codeVerifierBytes := make([]byte, 128)
		for i := range codeVerifierBytes {
			const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._~"
			codeVerifierBytes[i] = charset[rand.Intn(len(charset))]
		}
		codeVerifier := string(codeVerifierBytes)

		authURL := fmt.Sprintf(
			"https://www.dropbox.com/oauth2/authorize?client_id=%s&response_type=code&code_challenge=%s&code_challenge_method=%s&token_access_type=offline",
			clientID,
			codeVerifier, "plain",
		)

		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", authURL).Run()
		default:
			err = errors.New("no auto open for platform")
		}
		if err != nil {
			fmt.Println("Open this authentication URL in your browser:")
			fmt.Println(authURL)
			fmt.Println()
		}

		fmt.Println("Complete the OAuth flow in your browser to get an authorization code.")
		fmt.Println("Authorization Code:")
		fmt.Print("> ")
		var authCode string
		_, err = fmt.Scan(&authCode)
		if err != nil {
			return "", err
		}

		tokenReq = url.Values{
			"code":          []string{authCode},
			"grant_type":    []string{"authorization_code"},
			"code_verifier": []string{codeVerifier},
			"client_id":     []string{clientID},
		}
	} else {
		tokenReq = url.Values{
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{d.config.RefreshToken},
			"client_id":     []string{clientID},
		}
	}

	tokenRes, err := http.PostForm("https://www.dropbox.com/oauth2/token", tokenReq)
	if err != nil {
		return "", err
	}
	defer tokenRes.Body.Close()
	var tokenResFields struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(tokenRes.Body).Decode(&tokenResFields); err != nil {
		return "", err
	}

	d.accessToken = tokenResFields.AccessToken
	d.accessTokenExpiresAt = time.Now().Add(time.Duration(tokenResFields.ExpiresIn) * time.Second)

	if tokenResFields.RefreshToken != "" {
		d.config.RefreshToken = tokenResFields.RefreshToken
		configUpdated = true
	}

	return d.accessToken, nil
}
