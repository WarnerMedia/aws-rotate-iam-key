package main

import (
	"fmt"
)

var (
	// BuildVersion String
	Version string = ""
	// BuildTime populated by link
	BuildTime string = ""
	// RepoURL url of the repo
	RepoURL string = ""
)

// GetVersion - gets the version of the app
func GetVersion() string {
	return fmt.Sprintf("version %s \nbuilt %s \ngit repo = %s", Version, BuildTime, RepoURL)
}
