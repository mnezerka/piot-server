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
    MysqlDbHost string
    MysqlDbUsername string
    MysqlDbPassword string
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
        InfluxDbUri:        "",
        InfluxDbUsername:   "",
        InfluxDbPassword:   "",
        MysqlDbHost:        "",
        MysqlDbUsername:   "",
        MysqlDbPassword:   "",
        Params:         config.NewParameters(),
    }
    return o
}

