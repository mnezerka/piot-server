package service

import (
    "context"
    "github.com/op/go-logging"
    "piot-server/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

const VALUE_YES  = "yes"
const VALUE_NO = "no"

const TOPIC_UNIT = "unit"

const TOPIC_AVAILABLE = "available"

const TOPIC_NET = "net"
const TOPIC_IP = "net/ip"
const TOPIC_WIFI_SSID = "net/wifi/ssid"
const TOPIC_WIFI_STRENGTH = "net/wifi/strength"

const TOPIC_TEMP = "temperature"
const TOPIC_PRESSURE = "pressure"
const TOPIC_HUMIDITY = "humidity"


type IMqtt interface {
    PushThingData(ctx context.Context, thing *model.Thing, topic, value string) error
}

type Mqtt struct {
    MqttUri string
}

func (t *Mqtt) PushThingData(ctx context.Context, thing *model.Thing, topic, value string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Push thing data: %s", thing.Name)

    // post data to MQTT if device is enabled
    if thing.Enabled && thing.OrgId != primitive.NilObjectID {
        ctx.Value("log").(*logging.Logger).Debugf("TODO - write data to mqtt for enabled thing %v, topic %s, value %s", thing.Name, topic, value)
    }

    return nil
}
