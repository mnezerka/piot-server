package config

import (
    "time"
)

type Parameters struct {
    LogLevel string
    DOSInterval time.Duration
    JwtTokenExpiration time.Duration
    JwtPassword string
    DbUri string
    DbName string
    SmtpHost string
    SmtpPort int
    SmtpUser string
    SmtpPassword string
    MailFrom string
}

func NewParameters() *Parameters {
    p := &Parameters{
        LogLevel:       "INFO",
        DOSInterval:    1 * time.Second,
        JwtTokenExpiration: 5 * time.Hour,
        JwtPassword: "jwt-secret",
        DbUri: "",
        DbName: "",
        SmtpHost: "",
        SmtpPort: 25,
        SmtpUser: "",
        SmtpPassword: "",
        MailFrom: "",
    }
    return p
}
