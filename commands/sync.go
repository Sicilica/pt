package commands

import (
	"fmt"

	"github.com/sicilica/pt/types"
)

func init() {
	register("sync", commandSync, "Storage", "synchronizes local data with the cloud")
}

func commandSync(c types.CommandContext) error {
	return c.WithCloudSync(func(cloud types.CloudSyncInterface) error {
		localVer, err := cloud.GetLocalVersion()
		if err != nil {
			return err
		}

		remoteVer, err := cloud.GetRemoteVersion()
		if err != nil {
			return err
		}

		if localVer.Exists && (!remoteVer.Exists || localVer.Modified.After(remoteVer.Modified)) {
			fmt.Println("uploading local version to cloud")
			return cloud.Upload()
		}

		if remoteVer.Exists && (!localVer.Exists || remoteVer.Modified.After(localVer.Modified)) {
			fmt.Println("downloading updated version from cloud")
			return cloud.Download()
		}

		fmt.Println("already up-to-date")

		return nil
	})
}
