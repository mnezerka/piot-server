package context

type ContextOptions struct {
    MqttUri string
    MqttUsername string
    MqttPassword string
    DbUri string
    DbName string
    LogLevel string
}

func NewContextOptions() *ContextOptions {
    o := &ContextOptions{
        MqttUri:        "mock",
        MqttUsername:   "",
        MqttPassword:   "",
        DbUri:          "piot",
        DbName:         "piot",
        LogLevel:       "INFO",
    }
    return o
}

