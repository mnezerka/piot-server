package main_test

import (
    "github.com/op/go-logging"
    "piot-server/model"
)

type mysqlDbMockCall struct {
    Thing *model.Thing
    Value string
}

// implements IMysqlDb interface
type MysqlDbMock struct {
    Log *logging.Logger
    Calls []mysqlDbMockCall
}

func (db *MysqlDbMock) Open() error {
    return nil
}

func (db *MysqlDbMock) Close() {
}

func (db *MysqlDbMock) StoreMeasurement(thing *model.Thing, value string) {
    db.Log.Debugf("Mysqldb mock - store measurement, thing: %s, val: %s", thing.Name, value)
    db.Calls = append(db.Calls, mysqlDbMockCall{thing, value})
}

func (db *MysqlDbMock) StoreSwitchState(thing *model.Thing, value string) {
    db.Log.Debugf("Mysqldb mock - store switch state, thing: %s, val: %s", thing.Name, value)
    db.Calls = append(db.Calls, mysqlDbMockCall{thing, value})
}
