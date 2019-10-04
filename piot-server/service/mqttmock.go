package service

import (
    "context"
    "github.com/op/go-logging"
    "piot-server/model"
)

type call struct {
    Topic string
    Value string
    Thing *model.Thing
}

// implements IMqtt interface
type MqttMock struct {
    Calls []call
}

func (t *MqttMock) Connect(ctx context.Context) error {
    return nil
}

func (t *MqttMock) Disconnect(ctx context.Context) error {
    return nil
}

func (t *MqttMock) SetUsername(username string) {
}

func (t *MqttMock) SetPassword(password string) {
}

func (t *MqttMock) PushThingData(ctx context.Context, thing *model.Thing, topic, value string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Push thing data: %s, topic: %s, value: %s", thing.Name, topic, value)
    t.Calls = append(t.Calls, call{topic, value, thing})

    return nil
}
