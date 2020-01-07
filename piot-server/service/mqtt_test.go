package service_test

import (
    "fmt"
    "testing"
    "piot-server/service"
    "piot-server/test"
)

func TestMqttMsgNotSensor(t *testing.T) {

    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)

    mqtt := service.NewMqtt("uri")

    // send message to topic that is ignored
    mqtt.ProcessMessage(ctx, "xxx", "payload")

    // send message to not registered thing
    mqtt.ProcessMessage(ctx, "org/hello/x", "payload")
}

func TestMqttMsgSensor(t *testing.T) {
    const SENSOR = "sensor1"
    const ORG = "org1"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, SENSOR)
    orgId := test.CreateOrg(t, ctx, ORG)
    test.AddOrgThing(t, ctx, orgId, SENSOR)

    mqtt := service.NewMqtt("uri")

    // send unit message to registered thing
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/value/unit", ORG, SENSOR), "C")

    // send temperature message to registered thing
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/value/temperature", ORG, SENSOR), "23")

    // second round of calls to check proper functionality for high load

    // send unit message to registered thing (second call to check proper functionality
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/value/unit", ORG, SENSOR), "C")

    // send temperature message to registered thing
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/value/temperature", ORG, SENSOR), "23")

}
