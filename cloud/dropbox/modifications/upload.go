package modifications

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/file_properties"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

// FixedUpload is a modified version of the upload logic from the official SDK, but
// has ClientModified fixed to the right format.
//
// The original SDK is published under the MIT license, included here as "dropbox-LICENSE".
func FixedUpload(dbx *dropbox.Context, brokenArg *files.CommitInfo, content io.Reader) (res *files.FileMetadata, err error) {
	type fixedCommitInfo struct {
		Path string `json:"path"`
		Mode *files.WriteMode `json:"mode"`
		Autorename bool `json:"autorename"`
		ClientModified string `json:"client_modified,omitempty"`
		Mute bool `json:"mute"`
		PropertyGroups []*file_properties.PropertyGroup `json:"property_groups,omitempty"`
		StrictConflict bool `json:"strict_conflict"`
	}

	arg := &fixedCommitInfo{
		Path: brokenArg.Path,
		Mode: brokenArg.Mode,
		Autorename: brokenArg.Autorename,
		ClientModified: brokenArg.ClientModified.UTC().Format(time.RFC3339),
		Mute: brokenArg.Mute,
		PropertyGroups: brokenArg.PropertyGroups,
		StrictConflict: brokenArg.StrictConflict,
	}

	cli := dbx.Client

	dbx.Config.LogDebug("arg: %v", arg)
	b, err := json.Marshal(arg)
	if err != nil {
		return
	}

	headers := map[string]string{
		"Content-Type":    "application/octet-stream",
		"Dropbox-API-Arg": dropbox.HTTPHeaderSafeJSON(b),
	}
	if dbx.Config.AsMemberID != "" {
		headers["Dropbox-API-Select-User"] = dbx.Config.AsMemberID
	}

	req, err := (*dropbox.Context)(dbx).NewRequest("content", "upload", true, "files", "upload", headers, content)
	if err != nil {
		return
	}
	dbx.Config.LogInfo("req: %v", req)

	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	dbx.Config.LogInfo("resp: %v", resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	dbx.Config.LogDebug("body: %s", body)
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &res)
		if err != nil {
			return
		}

		return
	}
	if resp.StatusCode == http.StatusConflict {
		var apiError files.UploadAPIError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return
		}
		err = apiError
		return
	}
	err = auth.HandleCommonAuthErrors(dbx.Config, resp, body)
	if err != nil {
		return
	}
	err = dropbox.HandleCommonAPIErrors(dbx.Config, resp, body)
	return
}
