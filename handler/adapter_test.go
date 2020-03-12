package handler_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "github.com/mnezerka/go-piot/test"
    "piot-server/handler"
)

func getAdapter(t *testing.T) *handler.Adapter {
    log := test.GetLogger(t)
    db := test.GetDb(t)
    things := test.GetThings(t, log, db)
    mqtt := test.GetMqtt(t, log)
    pdevices := test.GetPiotDevices(t, log, things, mqtt)

    return handler.NewAdapter(log, pdevices)
}

/* GET method is not supported */
func TestForbiddenGet(t *testing.T) {
    req, err := http.NewRequest("GET", "/", strings.NewReader(""))
    test.Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 405)
}


/* Post data for device that is not registered  */
func TestPacketForUnknownThing(t *testing.T) {
    db := test.GetDb(t)
    test.CleanDb(t, db)

    deviceData := `
    {
        "device": "Device123",
        "readings": [
            {
                "address": "SensorXYZ",
                "t": 23
            }
        ]
    }`

    req, err := http.NewRequest("POST", "/", strings.NewReader(deviceData))
    test.Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}

/* Post data in short notation for device that is not registered */
func TestPacketShortNotationForUnknownThing(t *testing.T) {
    db := test.GetDb(t)
    test.CleanDb(t, db)

    deviceData := `
    {
        "d": "Device123",
        "r": [
            {
                "a": "SensorXYZ",
                "t": 23
            }
        ]
    }`

    req, err := http.NewRequest("POST", "/", strings.NewReader(deviceData))
    test.Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}
