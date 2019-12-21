package handler_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "piot-server/test"
    "piot-server/handler"
)

/* GET method is not supported */
func TestForbiddenGet(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    test.Ok(t, err)
    req = req.WithContext(test.CreateTestContext())

    rr := httptest.NewRecorder()
    handler := handler.Adapter{}
    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 405)
}


/* Post data for device that is not registered  */
func TestPacketForUnknownThing(t *testing.T) {
    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

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

    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()

    handler := handler.Adapter{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}

/* Post data in short notation for device that is not registered */
func TestPacketShortNotationForUnknownThing(t *testing.T) {
    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

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

    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()

    handler := handler.Adapter{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}
