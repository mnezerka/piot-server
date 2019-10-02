package handler_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "piot-server/test"
    "piot-server/handler"
)

func TestForbiddenGet(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    test.Ok(t, err)
    req = req.WithContext(test.CreateTestContext())

    rr := httptest.NewRecorder()
    handler := handler.Adapter{}
    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 405)
}

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

    // TODO: Check if defice is registered
}
