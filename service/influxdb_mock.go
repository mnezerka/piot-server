package service

import (
    "context"
    "github.com/op/go-logging"
    "piot-server/model"
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
    ctx.Value("log").(*logging.Logger).Debugf("Influxdb - push measurement, thing: %s, val: %s", thing.Name, value)

    db.Calls = append(db.Calls, influxDbMockCall{thing, value})
    //return nil
}
