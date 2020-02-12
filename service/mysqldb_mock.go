package service

import (
    "context"
    "github.com/op/go-logging"
    "piot-server/model"
)

type mysqlDbMockCall struct {
    Thing *model.Thing
    Value string
}

// implements IMysqlDb interface
type MysqlDbMock struct {
    Calls []mysqlDbMockCall
}

func (db *MysqlDbMock) Open(ctx context.Context) error {
    return nil
}

func (db *MysqlDbMock) Close(ctx context.Context) {
}

func (db *MysqlDbMock) StoreMeasurement(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Mysqldb mock - store measurement, thing: %s, val: %s", thing.Name, value)
    db.Calls = append(db.Calls, mysqlDbMockCall{thing, value})
}

func (db *MysqlDbMock) StoreSwitchState(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Mysqldb mock - store switch state, thing: %s, val: %s", thing.Name, value)
    db.Calls = append(db.Calls, mysqlDbMockCall{thing, value})
}
