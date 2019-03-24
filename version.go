package main

import (
	"fmt"
)

const (
	appName         = "astolfo"
	appVersionMajor = 0
	appVersionMinor = 1
	appRevision     = 0
)

func version() string {
	return fmt.Sprintf("%s version %d.%d.%d", appName, appVersionMajor,
		appVersionMinor, appRevision)
}
