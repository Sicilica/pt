package modifications

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/file_properties"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// FixedUpload is a modified version of the upload logic from the official SDK, but
// has ClientModified fixed to the right format.
//
// The original SDK is published under the MIT license, included here as "dropbox-LICENSE".
func FixedUpload(dbx *dropbox.Context, brokenArg *files.CommitInfo, content io.Reader) error {
	type fixedCommitInfo struct {
		Path           string                           `json:"path"`
		Mode           *files.WriteMode                 `json:"mode"`
		Autorename     bool                             `json:"autorename"`
		ClientModified string                           `json:"client_modified,omitempty"`
		Mute           bool                             `json:"mute"`
		PropertyGroups []*file_properties.PropertyGroup `json:"property_groups,omitempty"`
		StrictConflict bool                             `json:"strict_conflict"`
	}

	arg := &fixedCommitInfo{
		Path:           brokenArg.Path,
		Mode:           brokenArg.Mode,
		Autorename:     brokenArg.Autorename,
		ClientModified: brokenArg.ClientModified.UTC().Format(time.RFC3339),
		Mute:           brokenArg.Mute,
		PropertyGroups: brokenArg.PropertyGroups,
		StrictConflict: brokenArg.StrictConflict,
	}

	dbx.Config.LogDebug("arg: %v", arg)
	b, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Content-Type":    "application/octet-stream",
		"Dropbox-API-Arg": dropbox.HTTPHeaderSafeJSON(b),
	}
	if dbx.Config.AsMemberID != "" {
		headers["Dropbox-API-Select-User"] = dbx.Config.AsMemberID
	}

	req := dropbox.Request{
		Host:      "content",
		Namespace: "files",
		Route:     "upload",
		Auth:      "user",
		Style:     "upload",
		// TODO redundant?
		Arg:          arg,
		ExtraHeaders: headers,
	}
	_, _, err = dbx.Execute(req, content)
	if err != nil {
		var appErr files.UploadAPIError
		err = auth.ParseError(err, &appErr)
		if err == &appErr {
			err = appErr
		}
		return err
	}

	// NOTE We don't bother parsing the response
	return nil
}
