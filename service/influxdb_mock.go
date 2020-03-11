package service

import (
    "context"
    "github.com/op/go-logging"
    "github.com/mnezerka/go-piot/model"
)

type influxDbMockCall struct {
    Thing *model.Thing
    Value string
}

// implements IMqtt interface
type InfluxDbMock struct {
    Calls []influxDbMockCall
}

func (db *InfluxDbMock) PostMeasurement(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Influxdb - post measurement, thing: %s, val: %s", thing.Name, value)
    db.Calls = append(db.Calls, influxDbMockCall{thing, value})
}

func (db *InfluxDbMock) PostSwitchState(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Influxdb - post switch state, thing: %s, val: %s", thing.Name, value)
    db.Calls = append(db.Calls, influxDbMockCall{thing, value})
}
