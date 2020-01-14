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
    InfluxDbUri string
    InfluxDbUsername string
    InfluxDbPassword string
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
        InfluxDbUri:        "mock",
        InfluxDbUsername:   "",
        InfluxDbPassword:   "",
        Params:         config.NewParameters(),
    }
    return o
}

