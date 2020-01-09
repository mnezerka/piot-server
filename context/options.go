package context

import (
    "piot-server/config"
)

type ContextOptions struct {
    MqttUri string
    MqttUsername string
    MqttPassword string
    MqttClient string
    DbUri string
    DbName string
    Params *config.Parameters
}

func NewContextOptions() *ContextOptions {
    o := &ContextOptions{
        MqttUri:        "mock",
        MqttUsername:   "",
        MqttPassword:   "",
        MqttClient:     "",
        DbUri:          "piot",
        DbName:         "piot",
        Params:         config.NewParameters(),

    }
    return o
}

