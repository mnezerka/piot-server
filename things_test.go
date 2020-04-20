package main_test

import (
    "testing"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "piot-server"
)

func TestGetExistingThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    things := main.NewThings(GetDb(t), GetLogger(t))

    id := CreateThing(t, db, "thing1")

    thing, err := things.Get(id)
    Ok(t, err)
    Assert(t, thing.Name == "thing1", "Wrong thing name")
}

func TestGetUnknownThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    things := main.NewThings(GetDb(t), GetLogger(t))

    id := primitive.NewObjectID()

    _, err := things.Get(id)
    Assert(t, err != nil, "Thing shall not be found")
}

func TestFindUnknownThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    things := main.NewThings(GetDb(t), GetLogger(t))
    _, err := things.Find("xx")
    Assert(t, err != nil, "Thing shall not be found")
}

func TestFindExistingThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    CreateThing(t, db, "thing1")
    things := main.NewThings(GetDb(t), GetLogger(t))
    _, err := things.Find("thing1")
    Ok(t, err)
}

func TestRegisterThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    things := main.NewThings(GetDb(t), GetLogger(t))
    thing, err := things.RegisterPiot("thing1", "sensor")
    Ok(t, err)
    Equals(t, "thing1", thing.PiotId)
    Assert(t, thing.Name == "thing1", "Wrong thing name")
    Assert(t, thing.Type == "sensor", "Wrong thing type")
}

func TestSetParent(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)

    const THING_NAME_PARENT = "parent"
    id_parent := CreateThing(t, db, THING_NAME_PARENT)

    const THING_NAME_CHILD = "child"
    id_child := CreateThing(t, db, THING_NAME_CHILD)

    things := main.NewThings(GetDb(t), GetLogger(t))

    err := things.SetParent(id_child, id_parent)
    Ok(t, err)

    thing, err := things.Get(id_child)
    Ok(t, err)
    Equals(t, THING_NAME_CHILD, thing.Name)
    Equals(t, id_parent, thing.ParentId)
    /*test.test.Equals(t, "available", thing.AvailabilityTopic)
    test.Equals(t, "yes", thing.AvailabilityYes)
    test.Equals(t, "no", thing.AvailabilityNo)
    */
}

func TestTouchThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)

    const THING_NAME = "parent"
    id := CreateThing(t, db, THING_NAME)

    things := main.NewThings(GetDb(t), GetLogger(t))

    err := things.TouchThing(id)
    Ok(t, err)

    thing, err := things.Get(id)
    Ok(t, err)
    Equals(t, THING_NAME, thing.Name)
    // TODO check date
}

func TestSetAvailabilityAttributes(t *testing.T) {
    const THING_NAME = "thing2"
    db := GetDb(t)
    CleanDb(t, db)
    thingId := CreateThing(t, db, THING_NAME)
    things := main.NewThings(GetDb(t), GetLogger(t))
    err := things.SetAvailabilityTopic(thingId, "available")
    Ok(t, err)
    err = things.SetAvailabilityYesNo(thingId, "yes", "no")
    Ok(t, err)

    thing, err := things.Find(THING_NAME)
    Ok(t, err)
    Equals(t, THING_NAME, thing.Name)
    Equals(t, "available", thing.AvailabilityTopic)
    Equals(t, "yes", thing.AvailabilityYes)
    Equals(t, "no", thing.AvailabilityNo)
}

func TestSetLocationAttributes(t *testing.T) {
    const THING_NAME = "thing2"
    db := GetDb(t)
    CleanDb(t, db)
    thingId := CreateThing(t, db, THING_NAME)
    things := main.NewThings(GetDb(t), GetLogger(t))

    err := things.SetLocationMqttTopic(thingId, "loctopic")
    Ok(t, err)

    err = things.SetLocationMqttValues(thingId, "latval", "lngval", "satval", "tsval")
    Ok(t, err)

    thing, err := things.Find(THING_NAME)
    Ok(t, err)
    Equals(t, THING_NAME, thing.Name)
    Equals(t, "loctopic", thing.LocationMqttTopic)
    Equals(t, "latval", thing.LocationMqttLatValue)
    Equals(t, "lngval", thing.LocationMqttLngValue)
    Equals(t, "satval", thing.LocationMqttSatValue)
    Equals(t, "tsval", thing.LocationMqttTsValue)
}

func TestSetLocation(t *testing.T) {
    const THING_NAME = "thing2"
    db := GetDb(t)
    CleanDb(t, db)
    thingId := CreateThing(t, db, THING_NAME)
    things := main.NewThings(GetDb(t), GetLogger(t))

    err := things.SetLocation(thingId, 23.12, 56.33333, 4, 0)
    Ok(t, err)

    thing, err := things.Find(THING_NAME)
    Ok(t, err)
    Equals(t, THING_NAME, thing.Name)
    Equals(t, 23.12, thing.LocationLatitude)
    Equals(t, 56.33333, thing.LocationLongitude)
    Equals(t, int32(4), thing.LocationSatelites)

    // check that current value is overwritten by more recent measurement
    err = things.SetLocation(thingId, 1.1, 2.2, 1, 1000)
    Ok(t, err)

    thing, err = things.Find(THING_NAME)
    Ok(t, err)
    Equals(t, THING_NAME, thing.Name)
    Equals(t, 1.1, thing.LocationLatitude)

    // check that current value is not overwritten by old measurement
    err = things.SetLocation(thingId, 9.10, 10.11, 1, 900)
    Ok(t, err)

    thing, err = things.Find(THING_NAME)
    Ok(t, err)
    Equals(t, THING_NAME, thing.Name)
    Equals(t, 1.1, thing.LocationLatitude)
}

func TestSetSensorAttributes(t *testing.T) {
    const THING_NAME = "thing2"
    db := GetDb(t)
    CleanDb(t, db)
    thingId := CreateThing(t, db, THING_NAME)
    things := main.NewThings(GetDb(t), GetLogger(t))
    err := things.SetSensorMeasurementTopic(thingId, "value")
    Ok(t, err)

    err = things.SetSensorClass(thingId, "temperature")
    Ok(t, err)

    thing, err := things.Find(THING_NAME)
    Ok(t, err)
    Equals(t, THING_NAME, thing.Name)
    Equals(t, "value", thing.Sensor.MeasurementTopic)
    Equals(t, "temperature", thing.Sensor.Class)
}

func TestSetAlarm(t *testing.T) {
    const THING_NAME = "parent"
    db := GetDb(t)
    CleanDb(t, db)
    id := CreateThing(t, db, "thing")

    things := main.NewThings(GetDb(t), GetLogger(t))

    err := things.SetAlarm(id, true)
    Ok(t, err)

    thing, err := things.Get(id)
    Ok(t, err)
    Equals(t, true, thing.AlarmActive)
    // TODO check date

    err = things.SetAlarm(id, false)
    Ok(t, err)

    thing, err = things.Get(id)
    Ok(t, err)
    Equals(t, false, thing.AlarmActive)
    // TODO check date
}
