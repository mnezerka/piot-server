package service_test

import (
    "testing"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "piot-server/test"
    "piot-server/model"
    "piot-server/service"
)

// VALID packet + NEW device -> successful registration
func TestPacketDeviceReg(t *testing.T) {
    const DEVICE = "device01"
    const SENSOR = "SensorAddr"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

    // get instance of piot devices service
    s := ctx.Value("piotdevices").(*service.PiotDevices)

    // process packet for unknown device
    var packet model.PiotDevicePacket
    packet.Device = DEVICE

    var reading model.PiotSensorReading
    var temp float32 = 4.5
    reading.Address = SENSOR
    reading.Temperature = &temp
    packet.Readings = append(packet.Readings, reading)

    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // Check if device is registered
    db := ctx.Value("db").(*mongo.Database)
    var thing model.Thing
    err = db.Collection("things").FindOne(ctx, bson.M{"name": DEVICE}).Decode(&thing)
    test.Ok(t, err)
    test.Equals(t, DEVICE, thing.Name)
    test.Equals(t, model.THING_TYPE_DEVICE, thing.Type)
    test.Equals(t, "available", thing.AvailabilityTopic)

    var thing_sensor model.Thing
    err = db.Collection("things").FindOne(ctx, bson.M{"name": SENSOR}).Decode(&thing_sensor)
    test.Ok(t, err)
    test.Equals(t, SENSOR, thing_sensor.Name)
    test.Equals(t, model.THING_TYPE_SENSOR, thing_sensor.Type)
    test.Equals(t, "temperature", thing_sensor.Sensor.Class)
    test.Equals(t, "value", thing_sensor.Sensor.MeasurementTopic)

    // check correct assignment
    test.Equals(t, thing.Id, thing_sensor.ParentId)
}

// VALID packet + NEW device -> successful registration
// VALID packet + SENSOR reassigned -> change of parent
// This test simulates scenario where sensor is disconnected
// from one device and connected to another one
func TestPacketDeviceUpdateParent(t *testing.T) {
    const DEVICE = "device01"
    const DEVICE2 = "device02"
    const SENSOR = "SensorAddr"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

    // get instance of piot devices service
    s := ctx.Value("piotdevices").(*service.PiotDevices)

    // process packet for unknown device
    var packet model.PiotDevicePacket
    packet.Device = DEVICE

    var reading model.PiotSensorReading
    var temp float32 = 4.5
    reading.Address = SENSOR
    reading.Temperature = &temp
    packet.Readings = append(packet.Readings, reading)

    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // Check if device is registered
    db := ctx.Value("db").(*mongo.Database)
    var thing model.Thing
    err = db.Collection("things").FindOne(ctx, bson.M{"name": DEVICE}).Decode(&thing)
    test.Ok(t, err)
    test.Equals(t, DEVICE, thing.Name)
    test.Equals(t, model.THING_TYPE_DEVICE, thing.Type)
    test.Equals(t, "available", thing.AvailabilityTopic)

    var thing_sensor model.Thing
    err = db.Collection("things").FindOne(ctx, bson.M{"name": SENSOR}).Decode(&thing_sensor)
    test.Ok(t, err)
    test.Equals(t, SENSOR, thing_sensor.Name)
    test.Equals(t, model.THING_TYPE_SENSOR, thing_sensor.Type)
    test.Equals(t, "temperature", thing_sensor.Sensor.Class)
    test.Equals(t, "value", thing_sensor.Sensor.MeasurementTopic)

    // check correct assignment
    test.Equals(t, thing.Id, thing_sensor.ParentId)

    // assign sensor to new device
    packet.Device = DEVICE2
    err = s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // Check if second device is registered
    var thing2 model.Thing
    err = db.Collection("things").FindOne(ctx, bson.M{"name": DEVICE2}).Decode(&thing2)
    test.Ok(t, err)
    test.Equals(t, DEVICE2, thing2.Name)

    err = db.Collection("things").FindOne(ctx, bson.M{"name": SENSOR}).Decode(&thing_sensor)
    test.Ok(t, err)
    test.Equals(t, SENSOR, thing_sensor.Name)

    // check correct re-assignment
    test.Equals(t, thing2.Id, thing_sensor.ParentId)
}


// VALID packet + UNASSIGNED device -> no mqtt messages are published
func TestPacketDeviceDataUnassigned(t *testing.T) {

    const DEVICE = "device01"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

    // create unassigned thing
    test.CreateThing(t, ctx, DEVICE)

    // get instance of piot devices service
    s := ctx.Value("piotdevices").(*service.PiotDevices)

    // process packet for known device
    var packet model.PiotDevicePacket
    packet.Device = DEVICE

    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // check if mqtt was NOT called
    mqtt := ctx.Value("mqtt").(*service.MqttMock)
    test.Equals(t, 0, len(mqtt.Calls))
}

// VALID packet + ASSIGNED device -> mqtt messages are published
func TestPacketDeviceDataAssigned(t *testing.T) {

    const DEVICE = "device01"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

    // create and assign thing to org
    test.CreateThing(t, ctx, DEVICE)
    orgId := test.CreateOrg(t, ctx, "org1")
    test.AddOrgThing(t, ctx, orgId, DEVICE)

    // get instance of piot devices service
    s := ctx.Value("piotdevices").(*service.PiotDevices)

    // process packet for assigned device + provide wifi information
    var packet model.PiotDevicePacket
    packet.Device = DEVICE
    ssid := "SSID"
    packet.WifiSSID = &ssid

    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // check if mqtt was called
    mqtt := ctx.Value("mqtt").(*service.MqttMock)
    test.Equals(t, 2, len(mqtt.Calls))

    test.Equals(t, "available", mqtt.Calls[0].Topic)
    test.Equals(t, "yes", mqtt.Calls[0].Value)

    test.Equals(t, "net/wifi/ssid", mqtt.Calls[1].Topic)
    test.Equals(t, "SSID", mqtt.Calls[1].Value)
}

// VALID packet + UNASSIGNED device + TEMPERATURE -> no mqtt messages are published
func TestPacketDeviceReadingTempUnassigned(t *testing.T) {

    const DEVICE = "device01"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, DEVICE)

    // get instance of piot devices service
    s := ctx.Value("piotdevices").(*service.PiotDevices)

    // process packet for know device
    var temp float32 = 4.5
    var reading model.PiotSensorReading
    reading.Address = "SensorAddr"
    reading.Temperature = &temp

    var packet model.PiotDevicePacket
    packet.Device = DEVICE
    packet.Readings = append(packet.Readings, reading)

    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // check if mqtt was called
    mqtt := ctx.Value("mqtt").(*service.MqttMock)
    test.Equals(t, 0, len(mqtt.Calls))
}

// VALID packet + ASSIGNED device + TEMPERATURE -> mqtt messages are published
func TestPacketDeviceReadingTempAssigned(t *testing.T) {

    const DEVICE = "device01"
    const SENSOR = "SensorAddr"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, DEVICE)
    test.CreateThing(t, ctx, SENSOR)
    orgId := test.CreateOrg(t, ctx, "org1")
    test.AddOrgThing(t, ctx, orgId, DEVICE)
    test.AddOrgThing(t, ctx, orgId, SENSOR)

    // get instance of piot devices service
    s := ctx.Value("piotdevices").(*service.PiotDevices)

    // process packet for know device
    var temp float32 = 4.5
    //var pressure float32 = 20.8
    //var humidity float32 = 95.5
    var reading model.PiotSensorReading
    reading.Address = SENSOR
    reading.Temperature = &temp
    //reading.Pressure= &pressure
    //reading.Humidity = &humidity

    var packet model.PiotDevicePacket
    packet.Device = DEVICE
    packet.Readings = append(packet.Readings, reading)

    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // check if mqtt was called
    mqtt := ctx.Value("mqtt").(*service.MqttMock)
    test.Equals(t, 4, len(mqtt.Calls))

    test.Equals(t, "available", mqtt.Calls[0].Topic)
    test.Equals(t, "yes", mqtt.Calls[0].Value)
    test.Equals(t, DEVICE, mqtt.Calls[0].Thing.Name)

    test.Equals(t, "available", mqtt.Calls[1].Topic)
    test.Equals(t, "yes", mqtt.Calls[1].Value)
    test.Equals(t, "SensorAddr", mqtt.Calls[1].Thing.Name)

    test.Equals(t, "value", mqtt.Calls[2].Topic)
    test.Equals(t, "4.5", mqtt.Calls[2].Value)
    test.Equals(t, "SensorAddr", mqtt.Calls[2].Thing.Name)

    test.Equals(t, "value/unit", mqtt.Calls[3].Topic)
    test.Equals(t, "C", mqtt.Calls[3].Value)
    test.Equals(t, "SensorAddr", mqtt.Calls[3].Thing.Name)
}

// Test DOS (Denial Of Service) protection
func TestDOS(t *testing.T) {
    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

    // get instance of piot devices service
    s := ctx.Value("piotdevices").(*service.PiotDevices)

    var packet model.PiotDevicePacket

    // send first packet
    packet.Device = "device01"
    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // check that sending same packet in short time frame is not possible
    err = s.ProcessPacket(ctx, packet)
    test.Assert(t, err != nil, "DOS protection doesn't work")

    // check that sending packet for different device is possible
    packet.Device = "device02"
    err = s.ProcessPacket(ctx, packet)
    test.Ok(t, err)
}
