package main_test

import (
	"fmt"
	main "piot-server"

	"github.com/op/go-logging"
)

type influxDbMockCall struct {
	Thing *main.Thing
	Value string
}

// implements IMqtt interface
type InfluxDbMock struct {
	Log   *logging.Logger
	Calls []influxDbMockCall
}

func (db *InfluxDbMock) PostMeasurement(thing *main.Thing, value string) {
	db.Log.Debugf("Influxdb - post measurement, thing: %s, val: %s", thing.Name, value)
	db.Calls = append(db.Calls, influxDbMockCall{thing, value})
}

func (db *InfluxDbMock) PostSwitchState(thing *main.Thing, value string) {
	db.Log.Debugf("Influxdb - post switch state, thing: %s, val: %s", thing.Name, value)
	db.Calls = append(db.Calls, influxDbMockCall{thing, value})
}

func (db *InfluxDbMock) PostLocation(thing *main.Thing, lat, lng float64, sat, ts int32) {
	db.Log.Debugf("Influxdb - post location, thing: %s, val: %f %f %d %d", thing.Name, lat, lng, sat, ts)
	db.Calls = append(db.Calls, influxDbMockCall{thing, fmt.Sprintf("lat:%f-lng:%f-sat:%d-ts:%d", lat, lng, sat, ts)})
}

func (db *InfluxDbMock) PostBatteryLevel(thing *main.Thing, level int32) {
	db.Log.Debugf("Influxdb - post battery level, thing: %s, val: %d", thing.Name, level)
	db.Calls = append(db.Calls, influxDbMockCall{thing, fmt.Sprintf("level:%d", level)})
}
