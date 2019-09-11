package service

import (
    "github.com/op/go-logging"
    "os"
)

func NewLogger(log_format string, debug bool) *logging.Logger {
    backend := logging.NewLogBackend(os.Stderr, "", 0)
    format := logging.MustStringFormatter(log_format)
    backendFormatter := logging.NewBackendFormatter(backend, format)

    backendLeveled := logging.AddModuleLevel(backendFormatter)
    backendLeveled.SetLevel(logging.INFO, "")
    if debug {
        backendLeveled.SetLevel(logging.DEBUG, "")
    }

    logging.SetBackend(backendLeveled)
    logger := logging.MustGetLogger("server")
    return logger
}
