package config

import (
    "fmt"
    "runtime"
)

var (
    Version   = "1.0"
)

func VersionString() string {
    return fmt.Sprintf("%s, %s", Version, runtime.Version())
}
