package service_test

import (
    "testing"
    "piot-server/test"
    "piot-server/service"
    "piot-server/model"
)

// Push measurement for sensor
func TestPushMeasurementForSensor(t *testing.T) {
    const DEVICE = "device01"
    const SENSOR = "SensorAddr"

    ctx := test.CreateTestContext()

    // prepare data device + sensor assigned to org
    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, DEVICE)
    sensorId := test.CreateThing(t, ctx, SENSOR)
    orgId := test.CreateOrg(t, ctx, "org1")
    test.AddOrgThing(t, ctx, orgId, DEVICE)
    test.AddOrgThing(t, ctx, orgId, SENSOR)

    influxdb := service.NewInfluxDb("http://uri", "user", "pass")

    // get instance of piot devices service
    things := ctx.Value("things").(*service.Things)

    // get thing instance
    thing, err := things.Get(ctx, sensorId)
    test.Ok(t, err)

    // push measurement for thing
    influxdb.PostMeasurement(ctx, thing, "23")

    // check if http client was called
    httpClient := ctx.Value("httpclient").(*service.HttpClientMock)
    test.Equals(t, 1, len(httpClient.Calls))

    // check call parameters
    test.Equals(t, "http://uri/write?db=db", httpClient.Calls[0].Url)
    test.Equals(t, "sensor,id=" + sensorId.Hex() + ",name=SensorAddr,class=temperature value=23", httpClient.Calls[0].Body)
    test.Equals(t, "user", *httpClient.Calls[0].Username)
    test.Equals(t, "pass", *httpClient.Calls[0].Password)
}

// Push measurement for thing
func TestPushMeasurementForDevice(t *testing.T) {
    const DEVICE = "device01"

    ctx := test.CreateTestContext()

    // prepare data device + sensor assigned to org
    test.CleanDb(t, ctx)
    thingId := test.CreateThing(t, ctx, DEVICE)
    orgId := test.CreateOrg(t, ctx, "org1")
    test.AddOrgThing(t, ctx, orgId, DEVICE)

    influxdb := service.NewInfluxDb("http://uri", "user", "pass")

    // get instance of piot devices service
    things := ctx.Value("things").(*service.Things)

    // get thing instance
    thing, err := things.Get(ctx, thingId)
    test.Ok(t, err)

    // change type of the thing to device
    thing.Type = model.THING_TYPE_DEVICE

    // push measurement for thing
    influxdb.PostMeasurement(ctx, thing, "23")

    // check if http client was NOT called
    httpClient := ctx.Value("httpclient").(*service.HttpClientMock)
    test.Equals(t, 0, len(httpClient.Calls))
}
