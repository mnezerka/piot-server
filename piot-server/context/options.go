package context

import (
    "piot-server/config"
)

type ContextOptions struct {
    MqttUri string
    MqttUsername string
    MqttPassword string
    DbUri string
    DbName string
    Params *config.Parameters
}

func NewContextOptions() *ContextOptions {
    o := &ContextOptions{
        MqttUri:        "mock",
        MqttUsername:   "",
        MqttPassword:   "",
        DbUri:          "piot",
        DbName:         "piot",
        Params:         config.NewParameters(),

    }
    return o
}

