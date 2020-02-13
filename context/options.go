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
    MysqlDbName string
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
        MysqlDbName    :   "",
        Params:         config.NewParameters(),
    }
    return o
}

