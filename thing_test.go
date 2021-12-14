package main_test

import (
	main "piot-server"
	"testing"
)

func TestFetchThing(t *testing.T) {

	// prepare
	db := GetDb(t)
	CleanDb(t, db)
	id := CreateThing(t, db, "thing1")

	// test
	thing, err := main.NewThingFromDb(GetDb(t), GetLogger(t), id)
	Ok(t, err)
	Assert(t, thing.Name == "thing1", "Wrong thing name")
}

func TestFlushThing(t *testing.T) {

	// prepare
	db := GetDb(t)
	log := GetLogger(t)
	CleanDb(t, db)
	id := CreateThing(t, db, "thing1")

	// test
	thing, err := main.NewThingFromDb(db, log, id)
	Ok(t, err)

	thing.BatteryMqttTopic = "battopic"
	err = thing.Flush(db, log)
	Ok(t, err)

	thingVerify, err := main.NewThingFromDb(db, log, id)
	Ok(t, err)

	Assert(t, thingVerify.BatteryMqttTopic == "battopic", "Wrong battery topic")
}

func TestCreateThingWithInvalidType(t *testing.T) {

	// prepare
	db := GetDb(t)
	CleanDb(t, db)

	// test
	thing, err := main.NewThing(db, GetLogger(t), "thingname", "sometype")
	Fail(t, err)
	Assert(t, thing == nil, "nil not returned on error")
}

func TestCreateThingWithExistingName(t *testing.T) {

	// prepare
	db := GetDb(t)
	CleanDb(t, db)
	CreateThing(t, db, "thing1")

	// test
	thing, err := main.NewThing(db, GetLogger(t), "thing1", main.THING_TYPE_DEVICE)
	Fail(t, err)
	Assert(t, thing == nil, "nil not returned on error")
}

func TestCreateThing(t *testing.T) {

	// prepare
	db := GetDb(t)
	CleanDb(t, db)

	// test
	thing, err := main.NewThing(db, GetLogger(t), "thing1", main.THING_TYPE_DEVICE)
	Ok(t, err)
	Assert(t, thing.Name == "thing1", "Wrong thing name")
}
