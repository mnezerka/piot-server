package main_test

import (
	"context"
	"fmt"
	main "piot-server"
	"testing"
	"time"

	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getMqtt(t *testing.T, log *logging.Logger, db *mongo.Database, influxDb main.IInfluxDb, mysqlDb main.IMysqlDb) main.IMqtt {
	orgs := GetOrgs(t, log, db)
	things := GetThings(t, log, db)
	return main.NewMqtt("uri", log, things, orgs, influxDb, mysqlDb)
}

func TestMqttMsgNotSensor(t *testing.T) {

	db := GetDb(t)
	log := GetLogger(t)
	influxDb := GetInfluxDb(t, log)
	mysqlDb := GetMysqlDb(t, log)
	mqtt := getMqtt(t, log, db, influxDb, mysqlDb)

	CleanDb(t, db)

	// send message to topic that is ignored
	mqtt.ProcessMessage("xxx", "payload")

	// send message to not registered thing
	mqtt.ProcessMessage("org/hello/x", "payload")
}

func TestMqttThingTelemetry(t *testing.T) {
	const THING = "device1"
	const ORG = "org1"

	log := GetLogger(t)
	db := GetDb(t)
	influxDb := GetInfluxDb(t, log)
	mysqlDb := GetMysqlDb(t, log)
	mqtt := getMqtt(t, log, db, influxDb, mysqlDb)
	things := GetThings(t, log, db)

	CleanDb(t, db)
	thingId := CreateDevice(t, db, THING)
	SetThingTelemetryTopic(t, db, thingId, THING+"/"+"telemetry")
	orgId := CreateOrg(t, db, ORG)
	AddOrgThing(t, db, orgId, THING)

	// send telemetry message
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/telemetry", ORG, THING), "telemetry data")

	thing, err := things.Get(thingId)
	Ok(t, err)
	Equals(t, THING, thing.Name)
	Equals(t, "telemetry data", thing.Telemetry)
}

func TestMqttThingLocation(t *testing.T) {
	const THING = "device1"
	const THING2 = "device2"
	const ORG = "org1"

	log := GetLogger(t)
	db := GetDb(t)
	influxDb := GetInfluxDb(t, log)
	mysqlDb := GetMysqlDb(t, log)
	mqtt := getMqtt(t, log, db, influxDb, mysqlDb)
	things := GetThings(t, log, db)

	CleanDb(t, db)
	thingId := CreateDevice(t, db, THING)
	thing2Id := CreateDevice(t, db, THING2)
	orgId := CreateOrg(t, db, ORG)
	AddOrgThing(t, db, orgId, THING)
	AddOrgThing(t, db, orgId, THING2)

	// THING1 - timestamp mqtt extraction value SET, tracking is OFF
	SetThingLocationParams(t, db, thingId, THING+"/"+"loc", "lat", "lng", "sat", "ts", false)

	// THING2 - timestamp mqtt extraction value NOT SET, tracking is ON
	SetThingLocationParams(t, db, thing2Id, THING2+"/"+"loc", "lat", "lng", "sat", "", true)

	// THING1 send location message with timestamp -> timestamp is used
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/loc", ORG, THING), "{\"lat\": 123.234, \"lng\": 678.789, \"ts\": 456}")

	thing, err := things.Get(thingId)
	Ok(t, err)
	Equals(t, THING, thing.Name)
	Equals(t, 123.234, thing.LocationLatitude)
	Equals(t, 678.789, thing.LocationLongitude)
	Equals(t, int32(456), thing.LocationTs)

	// THING1 send location message without timestamp -> current time should be set
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/loc", ORG, THING), "{\"lat\": 123.234, \"lng\": 678.789}")

	thing, err = things.Get(thingId)
	Ok(t, err)
	Equals(t, THING, thing.Name)
	Equals(t, 123.234, thing.LocationLatitude)
	Equals(t, 678.789, thing.LocationLongitude)
	Equals(t, int32(time.Now().Unix()/60), int32(thing.LocationTs/60))

	// check no influxdb calls were initiated by previous steps
	Equals(t, 0, len(influxDb.Calls))

	// THING2 send location message with timestamp, -> current time should be set
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/loc", ORG, THING2), "{\"lat\": 211.1, \"lng\": 222.19, \"sat\": 4, \"ts\": 600}")
	thing, err = things.Get(thing2Id)
	Ok(t, err)
	Equals(t, THING2, thing.Name)
	Equals(t, 211.1, thing.LocationLatitude)
	Equals(t, 222.19, thing.LocationLongitude)
	Equals(t, int32(time.Now().Unix()/60), int32(thing.LocationTs/60))

	Equals(t, 1, len(influxDb.Calls))
	Equals(t, THING2, influxDb.Calls[0].Thing.Name)
	Contains(t, influxDb.Calls[0].Value, "lat:211.1")
	Contains(t, influxDb.Calls[0].Value, "lng:222.1")
	Contains(t, influxDb.Calls[0].Value, "sat:4")
}

// incoming sensor MQTT message for registered sensor
func TestMqttMsgSensor(t *testing.T) {
	const SENSOR = "sensor1"
	const ORG = "org1"

	log := GetLogger(t)
	db := GetDb(t)
	influxDb := GetInfluxDb(t, log)
	mysqlDb := GetMysqlDb(t, log)
	mqtt := getMqtt(t, log, db, influxDb, mysqlDb)

	CleanDb(t, db)
	sensorId := CreateThing(t, db, SENSOR)
	SetSensorMeasurementTopic(t, db, sensorId, SENSOR+"/"+"value")
	orgId := CreateOrg(t, db, ORG)
	AddOrgThing(t, db, orgId, SENSOR)

	// send unit message to registered thing
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/unit", ORG, SENSOR), "C")

	// send temperature message to registered thing
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/value", ORG, SENSOR), "23")

	// check if influxdb was called
	Equals(t, 1, len(influxDb.Calls))
	Equals(t, "23", influxDb.Calls[0].Value)
	Equals(t, SENSOR, influxDb.Calls[0].Thing.Name)

	// check if mysql was called
	Equals(t, 1, len(mysqlDb.Calls))
	Equals(t, "23", mysqlDb.Calls[0].Value)
	Equals(t, SENSOR, mysqlDb.Calls[0].Thing.Name)

	// second round of calls to check proper functionality for high load
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/unit", ORG, SENSOR), "C")
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/value", ORG, SENSOR), "23")
}

// this verifies that parsing json payloads works well
func TestMqttMsgSensorWithComplexValue(t *testing.T) {
	const SENSOR = "sensor1"
	const ORG = "org1"

	log := GetLogger(t)
	db := GetDb(t)
	influxDb := GetInfluxDb(t, log)
	mysqlDb := GetMysqlDb(t, log)
	mqtt := getMqtt(t, log, db, influxDb, mysqlDb)

	CleanDb(t, db)
	sensorId := CreateThing(t, db, SENSOR)
	SetSensorMeasurementTopic(t, db, sensorId, SENSOR+"/"+"value")
	orgId := CreateOrg(t, db, ORG)
	AddOrgThing(t, db, orgId, SENSOR)

	// modify sensor thing - set value template
	_, err := db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": sensorId}, bson.M{"$set": bson.M{"sensor.measurement_value": "temp"}})
	Ok(t, err)

	// send temperature message to registered thing
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/value", ORG, SENSOR), "{\"temp\": \"23\"}")

	// check if persistent storages were called
	Equals(t, 1, len(influxDb.Calls))
	Equals(t, 1, len(mysqlDb.Calls))
	Equals(t, "23", influxDb.Calls[0].Value)
	Equals(t, SENSOR, influxDb.Calls[0].Thing.Name)
	Equals(t, "23", mysqlDb.Calls[0].Value)
	Equals(t, SENSOR, mysqlDb.Calls[0].Thing.Name)

	// more complex structure
	_, err = db.Collection("things").UpdateOne(context.TODO(), bson.M{"_id": sensorId}, bson.M{"$set": bson.M{"sensor.measurement_value": "DS18B20.Temperature"}})
	Ok(t, err)

	payload := "{\"Time\":\"2020-01-24T22:52:58\",\"DS18B20\":{\"Id\":\"0416C18091FF\",\"Temperature\":23.0}"
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/value", ORG, SENSOR), payload)

	// check if persistent storages were called
	Equals(t, 2, len(influxDb.Calls))
	Equals(t, 2, len(mysqlDb.Calls))
	Equals(t, "23", influxDb.Calls[1].Value)
	Equals(t, SENSOR, influxDb.Calls[1].Thing.Name)
	Equals(t, "23", mysqlDb.Calls[1].Value)
	Equals(t, SENSOR, mysqlDb.Calls[1].Thing.Name)
}

// test for case when more sensors share same topic
func TestMqttMsgMultipleSensors(t *testing.T) {
	const SENSOR1 = "sensor1"
	const SENSOR2 = "sensor2"
	const ORG = "org1"

	log := GetLogger(t)
	db := GetDb(t)
	influxDb := GetInfluxDb(t, log)
	mysqlDb := GetMysqlDb(t, log)
	mqtt := getMqtt(t, log, db, influxDb, mysqlDb)

	CleanDb(t, db)
	sensor1Id := CreateThing(t, db, SENSOR1)
	SetSensorMeasurementTopic(t, db, sensor1Id, "xyz/value")
	sensor2Id := CreateThing(t, db, SENSOR2)
	SetSensorMeasurementTopic(t, db, sensor2Id, "xyz/value")

	orgId := CreateOrg(t, db, ORG)
	AddOrgThing(t, db, orgId, SENSOR1)
	AddOrgThing(t, db, orgId, SENSOR2)

	// send temperature message to registered thing
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/xyz/value", ORG), "23")

	// check if persistent storages were called
	Equals(t, 2, len(influxDb.Calls))
	Equals(t, "23", influxDb.Calls[0].Value)
	Equals(t, SENSOR1, influxDb.Calls[0].Thing.Name)
	Equals(t, "23", influxDb.Calls[1].Value)
	Equals(t, SENSOR2, influxDb.Calls[1].Thing.Name)

	Equals(t, 2, len(mysqlDb.Calls))
	Equals(t, "23", mysqlDb.Calls[0].Value)
	Equals(t, SENSOR1, mysqlDb.Calls[0].Thing.Name)
	Equals(t, "23", mysqlDb.Calls[1].Value)
	Equals(t, SENSOR2, mysqlDb.Calls[1].Thing.Name)
}

func TestMqttMsgSwitch(t *testing.T) {
	const THING = "THING1"
	const ORG = "org1"

	log := GetLogger(t)
	db := GetDb(t)
	influxDb := GetInfluxDb(t, log)
	mysqlDb := GetMysqlDb(t, log)
	mqtt := getMqtt(t, log, db, influxDb, mysqlDb)

	CleanDb(t, db)
	sensorId := CreateSwitch(t, db, THING)
	SetSwitchStateTopic(t, db, sensorId, THING+"/"+"state", "ON", "OFF")
	orgId := CreateOrg(t, db, ORG)
	AddOrgThing(t, db, orgId, THING)

	// send state change to ON
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/state", ORG, THING), "ON")

	// send state change to OFF
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/state", ORG, THING), "OFF")

	// check if mqtt was called
	Equals(t, 2, len(influxDb.Calls))

	Equals(t, "1", influxDb.Calls[0].Value)
	Equals(t, THING, influxDb.Calls[0].Thing.Name)

	Equals(t, "0", influxDb.Calls[1].Value)
	Equals(t, THING, influxDb.Calls[1].Thing.Name)
}

func TestMqttMsgBatteryLevel(t *testing.T) {

	/////////////////////////////////// prepare
	const THING = "THING1"
	const ORG = "org1"

	log := GetLogger(t)
	db := GetDb(t)
	influxDb := GetInfluxDb(t, log)
	mysqlDb := GetMysqlDb(t, log)
	mqtt := getMqtt(t, log, db, influxDb, mysqlDb)

	CleanDb(t, db)
	thing, err := main.NewThing(db, log, THING, main.THING_TYPE_DEVICE)
	Ok(t, err)

	// set battery level topic and level value
	thing.BatteryMqttTopic = THING + "/" + "state"
	thing.BatteryMqttLevelValue = "bat"
	thing.BatteryLevelTracking = true
	err = thing.Flush(db, log)
	Ok(t, err)

	orgId := CreateOrg(t, db, ORG)
	AddOrgThing(t, db, orgId, THING)

	/////////////////////////////////// test

	// send battery status
	mqtt.ProcessMessage(fmt.Sprintf("org/%s/%s/state", ORG, THING), "{\"bat\": 23}")

	/////////////////////////////////// check

	// fetch thing for verification
	thingVerify, err := main.NewThingFromDb(db, log, thing.Id)
	Ok(t, err)

	Equals(t, THING, thingVerify.Name)
	Equals(t, int32(23), thingVerify.BatteryLevel)

	// check if mqtt was called
	Equals(t, 1, len(influxDb.Calls))
	Equals(t, THING, influxDb.Calls[0].Thing.Name)
	Contains(t, influxDb.Calls[0].Value, "level:23")
}
