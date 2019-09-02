package main

import (
	"fmt"
	"runtime"
)

var (
	Version   = "1.0"
)

func versionString() string {
	return fmt.Sprintf("%s, %s", Version, runtime.Version())
}
