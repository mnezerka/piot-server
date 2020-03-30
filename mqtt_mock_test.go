package main_test

import (
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
    Log *logging.Logger
    Calls []call
}

func (t *MqttMock) Connect(subscribe bool) error {
    return nil
}

func (t *MqttMock) Disconnect() error {
    return nil
}

func (t *MqttMock) SetUsername(username string) {
}

func (t *MqttMock) SetPassword(password string) {
}

func (t *MqttMock) SetClient(id string) {
}

func (t *MqttMock) PushThingData(thing *model.Thing, topic, value string) (error) {
    t.Log.Debugf("Push thing data: %s, topic: %s, value: %s", thing.Name, topic, value)
    t.Calls = append(t.Calls, call{topic, value, thing})

    return nil
}

func (t *MqttMock) ProcessMessage(topic, payload string) {
}

