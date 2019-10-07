package service

import (
    "github.com/op/go-logging"
    "os"
)

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
