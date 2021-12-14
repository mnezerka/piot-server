package main_test

import (
	main "piot-server"

	"github.com/op/go-logging"
)

type mysqlDbMockCall struct {
	Thing *main.Thing
	Value string
}

// implements IMysqlDb interface
type MysqlDbMock struct {
	Log   *logging.Logger
	Calls []mysqlDbMockCall
}

func (db *MysqlDbMock) Open() error {
	return nil
}

func (db *MysqlDbMock) Close() {
}

func (db *MysqlDbMock) StoreMeasurement(thing *main.Thing, value string) {
	db.Log.Debugf("Mysqldb mock - store measurement, thing: %s, val: %s", thing.Name, value)
	db.Calls = append(db.Calls, mysqlDbMockCall{thing, value})
}

func (db *MysqlDbMock) StoreSwitchState(thing *main.Thing, value string) {
	db.Log.Debugf("Mysqldb mock - store switch state, thing: %s, val: %s", thing.Name, value)
	db.Calls = append(db.Calls, mysqlDbMockCall{thing, value})
}
