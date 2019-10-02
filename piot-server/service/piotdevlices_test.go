package service_test

import (
    "testing"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "piot-server/test"
    "piot-server/model"
    "piot-server/service"
)

func TestProcessPacketUnknownDevice(t *testing.T) {
    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

    s := ctx.Value("piotdevices").(*service.PiotDevices)

    // process packet for unknown device
    var packet model.PiotDevicePacket
    packet.Device = "device01"

    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // Check if defice is registered
    db := ctx.Value("db").(*mongo.Database)
    var thing model.Thing
    err = db.Collection("things").FindOne(ctx, bson.M{"name": packet.Device}).Decode(&thing)
    test.Ok(t, err)
    test.Equals(t, packet.Device, thing.Name)
}

func TestProcessPacketKnownDevice(t *testing.T) {

    const DEVICE = "device01"

    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, DEVICE)

    s := ctx.Value("piotdevices").(*service.PiotDevices)

    // process packet for know device
    var packet model.PiotDevicePacket
    packet.Device = DEVICE
    ssid := "SSID"
    packet.WifiSSID = &ssid

    err := s.ProcessPacket(ctx, packet)
    test.Ok(t, err)

    // check if mqtt was called
    mqtt := ctx.Value("mqtt").(*test.MqttMock)
    test.Equals(t, 2, len(mqtt.Calls))
    test.Equals(t, "available", mqtt.Calls[0].Topic)
    test.Equals(t, "yes", mqtt.Calls[0].Value)
    test.Equals(t, "net/wifi/ssid", mqtt.Calls[1].Topic)
    test.Equals(t, "SSID", mqtt.Calls[1].Value)
}


// Test Denial Of Service protection
func TestDOS(t *testing.T) {
    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

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
