package service_test

import (
    "fmt"
    "testing"
    "piot-server/service"
    "piot-server/test"
)

func TestMqttMsgNotSensor(t *testing.T) {
    const SENSOR = "sensor1"
    const ORG = "org1"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, SENSOR)
    orgId := test.CreateOrg(t, ctx, ORG)
    test.AddOrgThing(t, ctx, orgId, SENSOR)

    mqtt := service.NewMqtt("uri")

    // send message to topic that is ignored 
    mqtt.ProcessMessage(ctx, "xxx", "payload")

    // send message to not registered thing
    mqtt.ProcessMessage(ctx, "org/hello/x", "payload")

    // send message to registered thing
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s", ORG, SENSOR), "23")
}
