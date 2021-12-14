package main_test

import (
	"context"
	main "piot-server"
	"testing"

	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type services struct {
	log      *logging.Logger
	db       *mongo.Database
	orgs     *main.Orgs
	things   *main.Things
	influxDb main.IInfluxDb
	mysqlDb  main.IMysqlDb
	mqtt     *MqttMock
	pdevices *main.PiotDevices
}

func getServices(t *testing.T) *services {
	services := services{}
	services.log = GetLogger(t)
	services.db = GetDb(t)
	services.orgs = GetOrgs(t, services.log, services.db)
	services.things = GetThings(t, services.log, services.db)
	services.influxDb = GetInfluxDb(t, services.log)
	services.mysqlDb = GetMysqlDb(t, services.log)
	services.mqtt = GetMqtt(t, services.log)
	services.pdevices = GetPiotDevices(t, services.log, services.things, services.mqtt)

	return &services
}

// VALID packet + NEW device -> successful registration
func TestPacketDeviceReg(t *testing.T) {
	const DEVICE = "device01"
	const SENSOR = "SensorAddr"

	s := getServices(t)

	CleanDb(t, s.db)

	// process packet for unknown device
	var packet main.PiotDevicePacket
	packet.Device = DEVICE

	var reading main.PiotSensorReading
	var temp float32 = 4.5
	reading.Address = SENSOR
	reading.Temperature = &temp
	packet.Readings = append(packet.Readings, reading)

	err := s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// Check if device is registered
	var thing main.Thing
	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": DEVICE}).Decode(&thing)
	Ok(t, err)
	Equals(t, DEVICE, thing.Name)
	Equals(t, main.THING_TYPE_DEVICE, thing.Type)
	Equals(t, "available", thing.AvailabilityTopic)

	var thing_sensor main.Thing
	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": "T" + SENSOR}).Decode(&thing_sensor)
	Ok(t, err)
	Equals(t, "T"+SENSOR, thing_sensor.Name)
	Equals(t, main.THING_TYPE_SENSOR, thing_sensor.Type)
	Equals(t, "temperature", thing_sensor.Sensor.Class)
	Equals(t, "value", thing_sensor.Sensor.MeasurementTopic)

	// check correct assignment
	Equals(t, thing.Id, thing_sensor.ParentId)
}

// VALID packet with more measurements per 1 sensor +
// NEW device -> successful registration of device and all sensors
func TestPacketDeviceRegMultiple(t *testing.T) {
	const DEVICE = "device01"
	const SENSOR = "SensorAddr"

	s := getServices(t)

	CleanDb(t, s.db)

	// process packet for unknown device
	var packet main.PiotDevicePacket
	packet.Device = DEVICE

	var reading main.PiotSensorReading
	reading.Address = SENSOR
	var temp float32 = 4.5
	reading.Temperature = &temp

	var press float32 = 900
	reading.Pressure = &press

	var hum float32 = 20
	reading.Humidity = &hum

	packet.Readings = append(packet.Readings, reading)

	err := s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// Check if device is registered
	var thing main.Thing
	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": DEVICE}).Decode(&thing)
	Ok(t, err)
	Equals(t, DEVICE, thing.Name)
	Equals(t, main.THING_TYPE_DEVICE, thing.Type)
	Equals(t, "available", thing.AvailabilityTopic)

	var thing_sensor main.Thing
	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": "T" + SENSOR}).Decode(&thing_sensor)
	Ok(t, err)
	Equals(t, "T"+SENSOR, thing_sensor.Name)
	Equals(t, main.THING_TYPE_SENSOR, thing_sensor.Type)
	Equals(t, "temperature", thing_sensor.Sensor.Class)
	Equals(t, "value", thing_sensor.Sensor.MeasurementTopic)

	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": "P" + SENSOR}).Decode(&thing_sensor)
	Ok(t, err)
	Equals(t, "P"+SENSOR, thing_sensor.Name)
	Equals(t, main.THING_TYPE_SENSOR, thing_sensor.Type)
	Equals(t, "pressure", thing_sensor.Sensor.Class)
	Equals(t, "value", thing_sensor.Sensor.MeasurementTopic)

	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": "H" + SENSOR}).Decode(&thing_sensor)
	Ok(t, err)
	Equals(t, "H"+SENSOR, thing_sensor.Name)
	Equals(t, main.THING_TYPE_SENSOR, thing_sensor.Type)
	Equals(t, "humidity", thing_sensor.Sensor.Class)
	Equals(t, "value", thing_sensor.Sensor.MeasurementTopic)

	// check correct assignment
	Equals(t, thing.Id, thing_sensor.ParentId)
}

// VALID packet + NEW device -> successful registration
// VALID packet + SENSOR reassigned -> change of parent
// This test simulates scenario where sensor is disconnected
// from one device and connected to another one
func TestPacketDeviceUpdateParent(t *testing.T) {
	const DEVICE = "device01"
	const DEVICE2 = "device02"
	const SENSOR = "SensorAddr"

	s := getServices(t)

	CleanDb(t, s.db)

	// process packet for unknown device
	var packet main.PiotDevicePacket
	packet.Device = DEVICE

	var reading main.PiotSensorReading
	var temp float32 = 4.5
	reading.Address = SENSOR
	reading.Temperature = &temp
	packet.Readings = append(packet.Readings, reading)

	err := s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// Check if device is registered
	var thing main.Thing
	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": DEVICE}).Decode(&thing)
	Ok(t, err)
	Equals(t, DEVICE, thing.Name)
	Equals(t, main.THING_TYPE_DEVICE, thing.Type)
	Equals(t, "available", thing.AvailabilityTopic)

	var thing_sensor main.Thing
	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": "T" + SENSOR}).Decode(&thing_sensor)
	Ok(t, err)
	Equals(t, "T"+SENSOR, thing_sensor.Name)
	Equals(t, main.THING_TYPE_SENSOR, thing_sensor.Type)
	Equals(t, "temperature", thing_sensor.Sensor.Class)
	Equals(t, "value", thing_sensor.Sensor.MeasurementTopic)

	// check correct assignment
	Equals(t, thing.Id, thing_sensor.ParentId)

	// assign sensor to new device
	packet.Device = DEVICE2
	err = s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// Check if second device is registered
	var thing2 main.Thing
	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": DEVICE2}).Decode(&thing2)
	Ok(t, err)
	Equals(t, DEVICE2, thing2.Name)

	err = s.db.Collection("things").FindOne(context.TODO(), bson.M{"name": "T" + SENSOR}).Decode(&thing_sensor)
	Ok(t, err)
	Equals(t, "T"+SENSOR, thing_sensor.Name)

	// check correct re-assignment
	Equals(t, thing2.Id, thing_sensor.ParentId)
}

// VALID packet + UNASSIGNED device -> no mqtt messages are published
func TestPacketDeviceDataUnassigned(t *testing.T) {

	const DEVICE = "device01"

	s := getServices(t)

	CleanDb(t, s.db)

	// create unassigned thing
	CreateThing(t, s.db, DEVICE)

	// process packet for known device
	var packet main.PiotDevicePacket
	packet.Device = DEVICE

	err := s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// check if mqtt was NOT called
	Equals(t, 0, len(s.mqtt.Calls))
}

// VALID packet + ASSIGNED device -> mqtt messages are published
func TestPacketDeviceDataAssigned(t *testing.T) {

	const DEVICE = "device01"

	s := getServices(t)

	CleanDb(t, s.db)

	// create and assign thing to org
	CreateThing(t, s.db, DEVICE)
	orgId := CreateOrg(t, s.db, "org1")
	AddOrgThing(t, s.db, orgId, DEVICE)

	// process packet for assigned device + provide wifi information
	var packet main.PiotDevicePacket
	packet.Device = DEVICE
	ssid := "SSID"
	packet.WifiSSID = &ssid

	err := s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// check if mqtt was called
	Equals(t, 2, len(s.mqtt.Calls))

	Equals(t, "available", s.mqtt.Calls[0].Topic)
	Equals(t, "yes", s.mqtt.Calls[0].Value)

	Equals(t, "net/wifi/ssid", s.mqtt.Calls[1].Topic)
	Equals(t, "SSID", s.mqtt.Calls[1].Value)
}

// VALID packet + UNASSIGNED device + TEMPERATURE -> no mqtt messages are published
func TestPacketDeviceReadingTempUnassigned(t *testing.T) {

	const DEVICE = "device01"

	s := getServices(t)

	CleanDb(t, s.db)
	CreateThing(t, s.db, DEVICE)

	// process packet for know device
	var temp float32 = 4.5
	var reading main.PiotSensorReading
	reading.Address = "SensorAddr"
	reading.Temperature = &temp

	var packet main.PiotDevicePacket
	packet.Device = DEVICE
	packet.Readings = append(packet.Readings, reading)

	err := s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// check if mqtt was called
	Equals(t, 0, len(s.mqtt.Calls))
}

// VALID packet + ASSIGNED device + TEMPERATURE -> mqtt messages are published
func TestPacketDeviceReadingTempAssigned(t *testing.T) {

	const DEVICE = "device01"
	const SENSOR = "SensorAddr"

	s := getServices(t)

	CleanDb(t, s.db)
	CreateThing(t, s.db, DEVICE)
	CreateThing(t, s.db, "T"+SENSOR) // SENSOR is registered for temperature
	orgId := CreateOrg(t, s.db, "org1")
	AddOrgThing(t, s.db, orgId, DEVICE)
	AddOrgThing(t, s.db, orgId, "T"+SENSOR) // SENSOR is registered for temperature

	// process packet for know device
	var temp float32 = 4.5
	//var pressure float32 = 20.8
	//var humidity float32 = 95.5
	var reading main.PiotSensorReading
	reading.Address = SENSOR
	reading.Temperature = &temp
	//reading.Pressure= &pressure
	//reading.Humidity = &humidity

	var packet main.PiotDevicePacket
	packet.Device = DEVICE
	packet.Readings = append(packet.Readings, reading)

	err := s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// check if mqtt was called
	Equals(t, 4, len(s.mqtt.Calls))

	Equals(t, "available", s.mqtt.Calls[0].Topic)
	Equals(t, "yes", s.mqtt.Calls[0].Value)
	Equals(t, DEVICE, s.mqtt.Calls[0].Thing.Name)

	Equals(t, "available", s.mqtt.Calls[1].Topic)
	Equals(t, "yes", s.mqtt.Calls[1].Value)
	Equals(t, "TSensorAddr", s.mqtt.Calls[1].Thing.Name)

	Equals(t, "value", s.mqtt.Calls[2].Topic)
	Equals(t, "4.5", s.mqtt.Calls[2].Value)
	Equals(t, "TSensorAddr", s.mqtt.Calls[2].Thing.Name)

	Equals(t, "value/unit", s.mqtt.Calls[3].Topic)
	Equals(t, "C", s.mqtt.Calls[3].Value)
	Equals(t, "TSensorAddr", s.mqtt.Calls[3].Thing.Name)
}

// Test DOS (Denial Of Service) protection
func TestDOS(t *testing.T) {

	s := getServices(t)

	CleanDb(t, s.db)

	var packet main.PiotDevicePacket

	// send first packet
	packet.Device = "device01"
	err := s.pdevices.ProcessPacket(packet)
	Ok(t, err)

	// check that sending same packet in short time frame is not possible
	err = s.pdevices.ProcessPacket(packet)
	Assert(t, err != nil, "DOS protection doesn't work")

	// check that sending packet for different device is possible
	packet.Device = "device02"
	err = s.pdevices.ProcessPacket(packet)
	Ok(t, err)
}
