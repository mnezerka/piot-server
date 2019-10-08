package config

import (
    "time"
)

type Parameters struct {
    LogLevel string
    DOSInterval time.Duration
    JwtTokenExpiration time.Duration
    JwtPassword string
}

func NewParameters() *Parameters {
    p := &Parameters{
        LogLevel:       "INFO",
        DOSInterval:    1 * time.Second,
        JwtTokenExpiration: 5 * time.Hour,
        JwtPassword: "jwt-secret",
    }
    return p
}

