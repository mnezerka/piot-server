package service_test

import (
    "testing"
    "piot-server/service"
    "piot-server/test"
    "piot-server/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetExistingThing(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    things := service.Things{}

    id := test.CreateThing(t, ctx, "thing1")

    thing, err := things.Get(ctx, id)
    test.Ok(t, err)
    test.Assert(t, thing.Name == "thing1", "Wrong thing name")
}

func TestGetUnknownThing(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    things := service.Things{}

    id := primitive.NewObjectID()

    _, err := things.Get(ctx, id)
    test.Assert(t, err != nil, "Thing shall not be found")
}

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
    thing, err := things.RegisterPiot(ctx, "thing1", "sensor")
    test.Ok(t, err)
    test.Equals(t, "thing1", thing.PiotId)
    test.Assert(t, thing.Name == "thing1", "Wrong thing name")
    test.Assert(t, thing.Type == "sensor", "Wrong thing type")
}

func TestSetParent(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)

    const THING_NAME_PARENT = "parent"
    id_parent := test.CreateThing(t, ctx, THING_NAME_PARENT)

    const THING_NAME_CHILD = "child"
    id_child := test.CreateThing(t, ctx, THING_NAME_CHILD)

    things := service.Things{}

    err := things.SetParent(ctx, id_child, id_parent)
    test.Ok(t, err)

    thing, err := things.Get(ctx, id_child)
    test.Ok(t, err)
    test.Equals(t, THING_NAME_CHILD, thing.Name)
    test.Equals(t, id_parent, thing.ParentId)
    /*test.Equals(t, "available", thing.AvailabilityTopic)
    test.Equals(t, "yes", thing.AvailabilityYes)
    test.Equals(t, "no", thing.AvailabilityNo)
    */
}

func TestTouchThing(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)

    const THING_NAME = "parent"
    id := test.CreateThing(t, ctx, THING_NAME)

    things := service.Things{}

    err := things.TouchThing(ctx, id)
    test.Ok(t, err)

    thing, err := things.Get(ctx, id)
    test.Ok(t, err)
    test.Equals(t, THING_NAME, thing.Name)
    // TODO check date
}

func TestSetAvailabilityAttributes(t *testing.T) {
    const THING_NAME = "thing2"
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    thingId := test.CreateThing(t, ctx, THING_NAME)
    things := service.Things{}
    err := things.SetAvailabilityTopic(ctx, thingId, "available")
    test.Ok(t, err)
    err = things.SetAvailabilityYesNo(ctx, thingId, "yes", "no")
    test.Ok(t, err)

    thing, err := things.Find(ctx, THING_NAME)
    test.Ok(t, err)
    test.Equals(t, THING_NAME, thing.Name)
    test.Equals(t, "available", thing.AvailabilityTopic)
    test.Equals(t, "yes", thing.AvailabilityYes)
    test.Equals(t, "no", thing.AvailabilityNo)
}

func TestSetLocation(t *testing.T) {
    const THING_NAME = "thing2"
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    thingId := test.CreateThing(t, ctx, THING_NAME)
    things := service.Things{}

    loc := model.LocationData{23.12, 56.33333};

    err := things.SetLocation(ctx, thingId, loc)
    test.Ok(t, err)

    thing, err := things.Find(ctx, THING_NAME)
    test.Ok(t, err)
    test.Equals(t, THING_NAME, thing.Name)
    test.Equals(t, 23.12, thing.Location.Latitude)
    test.Equals(t, 56.33333, thing.Location.Longitude)
}

func TestSetSensorAttributes(t *testing.T) {
    const THING_NAME = "thing2"
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    thingId := test.CreateThing(t, ctx, THING_NAME)
    things := service.Things{}
    err := things.SetSensorMeasurementTopic(ctx, thingId, "value")
    test.Ok(t, err)
    /*
    err = things.SetSensorClass(ctx, THING_NAME, "temperature")
    test.Ok(t, err)

    thing, err := things.Find(ctx, THING_NAME)
    test.Ok(t, err)
    test.Equals(t, THING_NAME, thing.Name)
    test.Equals(t, "value", thing.Sensor.MeasurementTopic)
    test.Equals(t, "temperature", thing.Sensor.Class)
    */
}
