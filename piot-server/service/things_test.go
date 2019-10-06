package service_test

import (
    "testing"
    "piot-server/service"
    "piot-server/test"
)

func TestFindUnknownThing(t *testing.T) {
    ctx := test.CreateTestContext()
    things := service.Things{}
    _, err := things.Find(ctx, "xx")
    test.Assert(t, err != nil, "Thing shall not be found")
}

func TestFindExistingThing(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, "thing1")
    things := service.Things{}
    _, err := things.Find(ctx, "thing1")
    test.Ok(t, err)
}

func TestRegisterThing(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    things := service.Things{}
    thing, err := things.Register(ctx, "thing1", "sensor")
    test.Ok(t, err)
    test.Assert(t, thing.Name == "thing1", "Wrong thing name")
    test.Assert(t, thing.Type == "sensor", "Wrong thing type")
}

func TestSetAvailabilityAttributes(t *testing.T) {
    const THING_NAME = "thing2"
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, THING_NAME)
    things := service.Things{}
    err := things.SetAvailabilityTopic(ctx, THING_NAME, "available")
    test.Ok(t, err)
    err = things.SetAvailabilityYesNo(ctx, THING_NAME, "yes", "no")
    test.Ok(t, err)

    thing, err := things.Find(ctx, THING_NAME)
    test.Ok(t, err)
    test.Equals(t, THING_NAME, thing.Name)
    test.Equals(t, "available", thing.AvailabilityTopic)
    test.Equals(t, "yes", thing.AvailabilityYes)
    test.Equals(t, "no", thing.AvailabilityNo)
}

func TestSetSensorAttributes(t *testing.T) {
    const THING_NAME = "thing2"
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, THING_NAME)
    things := service.Things{}
    err := things.SetSensorMeasurementTopic(ctx, THING_NAME, "value")
    test.Ok(t, err)
    err = things.SetSensorClass(ctx, THING_NAME, "temperature")
    test.Ok(t, err)

    thing, err := things.Find(ctx, THING_NAME)
    test.Ok(t, err)
    test.Equals(t, THING_NAME, thing.Name)
    test.Equals(t, "value", thing.Sensor.MeasurementTopic)
    test.Equals(t, "temperature", thing.Sensor.Class)
}
