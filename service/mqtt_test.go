package service_test

import (
    "fmt"
    "testing"
    "piot-server/service"
    "piot-server/test"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
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

func TestMqttThingTelemetry(t *testing.T) {
    const THING = "device1"
    const ORG = "org1"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    thingId := test.CreateDevice(t, ctx, THING)
    test.SetThingTelemetryTopic(t, ctx, thingId, THING + "/" + "telemetry")
    orgId := test.CreateOrg(t, ctx, ORG)
    test.AddOrgThing(t, ctx, orgId, THING)

    mqtt := service.NewMqtt("uri")

    // send telemetry message
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/telemetry", ORG, THING), "telemetry data")

    things := service.Things{}
    thing, err := things.Get(ctx, thingId)
    test.Ok(t, err)
    test.Equals(t, THING, thing.Name)
    test.Equals(t, "telemetry data", thing.Telemetry)
}

func TestMqttMsgSensor(t *testing.T) {
    const SENSOR = "sensor1"
    const ORG = "org1"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    sensorId := test.CreateThing(t, ctx, SENSOR)
    test.SetSensorMeasurementTopic(t, ctx, sensorId, SENSOR + "/" + "value")
    orgId := test.CreateOrg(t, ctx, ORG)
    test.AddOrgThing(t, ctx, orgId, SENSOR)

    mqtt := service.NewMqtt("uri")
    influxDb := ctx.Value("influxdb").(*service.InfluxDbMock)

    // send unit message to registered thing
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/unit", ORG, SENSOR), "C")

    // send temperature message to registered thing
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/value", ORG, SENSOR), "23")

    // check if mqtt was called
    test.Equals(t, 1, len(influxDb.Calls))

    test.Equals(t, "23", influxDb.Calls[0].Value)
    test.Equals(t, SENSOR, influxDb.Calls[0].Thing.Name)

    // second round of calls to check proper functionality for high load
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/unit", ORG, SENSOR), "C")
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/value", ORG, SENSOR), "23")
}

// this verifies that parsing json payloads works well
func TestMqttMsgSensorWithComplexValue(t *testing.T) {
    const SENSOR = "sensor1"
    const ORG = "org1"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    sensorId := test.CreateThing(t, ctx, SENSOR)
    test.SetSensorMeasurementTopic(t, ctx, sensorId, SENSOR + "/" + "value")
    orgId := test.CreateOrg(t, ctx, ORG)
    test.AddOrgThing(t, ctx, orgId, SENSOR)

    // modify sensor thing - set value template
    db := ctx.Value("db").(*mongo.Database)
    _, err := db.Collection("things").UpdateOne(ctx, bson.M{"_id": sensorId}, bson.M{"$set": bson.M{"sensor.measurement_value": "temp"}})
    test.Ok(t, err)

    mqtt := service.NewMqtt("uri")
    influxDb := ctx.Value("influxdb").(*service.InfluxDbMock)

    // send temperature message to registered thing
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/value", ORG, SENSOR), "{\"temp\": \"23\"}")

    // check if mqtt was called
    test.Equals(t, 1, len(influxDb.Calls))

    test.Equals(t, "23", influxDb.Calls[0].Value)
    test.Equals(t, SENSOR, influxDb.Calls[0].Thing.Name)

    // more complex structure
    _, err = db.Collection("things").UpdateOne(ctx, bson.M{"_id": sensorId}, bson.M{"$set": bson.M{"sensor.measurement_value": "DS18B20.Temperature"}})
    test.Ok(t, err)

    payload := "{\"Time\":\"2020-01-24T22:52:58\",\"DS18B20\":{\"Id\":\"0416C18091FF\",\"Temperature\":23.0}"
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/%s/value", ORG, SENSOR), payload)

    // check if mqtt was called
    test.Equals(t, 2, len(influxDb.Calls))

    test.Equals(t, "23", influxDb.Calls[1].Value)
    test.Equals(t, SENSOR, influxDb.Calls[1].Thing.Name)
}

// test for case when more sensors sue same topic
func TestMqttMsgMultipleSensors(t *testing.T) {
    const SENSOR1 = "sensor1"
    const SENSOR2 = "sensor2"
    const ORG = "org1"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    sensor1Id := test.CreateThing(t, ctx, SENSOR1)
    test.SetSensorMeasurementTopic(t, ctx, sensor1Id, "xyz/value")
    sensor2Id := test.CreateThing(t, ctx, SENSOR2)
    test.SetSensorMeasurementTopic(t, ctx, sensor2Id, "xyz/value")

    orgId := test.CreateOrg(t, ctx, ORG)
    test.AddOrgThing(t, ctx, orgId, SENSOR1)
    test.AddOrgThing(t, ctx, orgId, SENSOR2)

    mqtt := service.NewMqtt("uri")
    influxDb := ctx.Value("influxdb").(*service.InfluxDbMock)

    // send temperature message to registered thing
    mqtt.ProcessMessage(ctx, fmt.Sprintf("org/%s/xyz/value", ORG), "23")

    // check if influxdb was called
    test.Equals(t, 2, len(influxDb.Calls))

    test.Equals(t, "23", influxDb.Calls[0].Value)
    test.Equals(t, SENSOR1, influxDb.Calls[0].Thing.Name)

    test.Equals(t, "23", influxDb.Calls[1].Value)
    test.Equals(t, SENSOR2, influxDb.Calls[1].Thing.Name)
}


