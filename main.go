package main

import (
	"fmt"
	"os"

	"github.com/sicilica/pt/app"
	"github.com/sicilica/pt/cloud/dropbox"
	"github.com/sicilica/pt/storage/sqlite3"
)

func main() {
	err := func() error {
		runtime := app.Runtime{
			NewStorageProvider:   sqlite3.New,
			NewCloudSyncProvider: dropbox.New,
		}
		defer runtime.Close()

		return runtime.ExecCommand(os.Args[1:])
	}()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
