package main

import (
	"fmt"
	"os"

	"github.com/foxcpp/tagdb/storage"
)

func userHomeDir() (home string) {
	home = os.Getenv("HOME")
	if len(home) == 0 {
		// Windows...
		homeDrive := os.Getenv("HOMEDRIVE")
		homePath := os.Getenv("HOMEPATH")
		if len(homeDrive) == 0 || len(homePath) == 0 {
			home = os.Getenv("USERPROFILE")
		} else {
			home = homeDrive + homePath
		}
	}
	return
}

func getDB() (*storage.S, error) {
	path := os.Getenv("TAGDB")
	if len(path) == 0 {
		path = userHomeDir() + string(os.PathSeparator) + ".tag.db"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Creating new tags database at", path+"...")
	}

	return storage.Open(path)
}
