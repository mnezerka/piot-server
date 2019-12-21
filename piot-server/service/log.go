package service

import (
    "github.com/op/go-logging"
    "os"
)

// NewLogger creates instance of logger that should be used
// in all server handlers and routines. The idea is to have
// unified style of logging - logger is configured only once
// and at one place
func NewLogger(log_format string, level string) (*logging.Logger, error) {
    backend := logging.NewLogBackend(os.Stderr, "", 0)
    format := logging.MustStringFormatter(log_format)
    backendFormatter := logging.NewBackendFormatter(backend, format)

    backendLeveled := logging.AddModuleLevel(backendFormatter)
    logLevel, err := logging.LogLevel(level)
    if err != nil {
        return nil, err
    }
    backendLeveled.SetLevel(logLevel, "")

    logging.SetBackend(backendLeveled)
    logger := logging.MustGetLogger("server")
    return logger, nil
}
