package main_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "piot-server"
)

func getAdapter(t *testing.T) *main.Adapter {
    log := GetLogger(t)
    db := GetDb(t)
    things := GetThings(t, log, db)
    mqtt := GetMqtt(t, log)
    pdevices := GetPiotDevices(t, log, things, mqtt)

    return main.NewAdapter(log, pdevices)
}

/* GET method is not supported */
func TestForbiddenGet(t *testing.T) {
    req, err := http.NewRequest("GET", "/", strings.NewReader(""))
    Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 405)
}


/* Post data for device that is not registered  */
func TestPacketForUnknownThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)

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
    Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}

/* Post data in short notation for device that is not registered */
func TestPacketShortNotationForUnknownThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)

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
    Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}
