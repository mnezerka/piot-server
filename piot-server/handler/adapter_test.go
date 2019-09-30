package handler_test

import (
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"
    "piot-server/test"
    "piot-server/handler"
    piotcontext "piot-server/context"
)

func TestForbiddenGet(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    test.Ok(t, err)
    req = req.WithContext(piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test"))

    rr := httptest.NewRecorder()
    handler := handler.Adapter{}
    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 405)
}

func TestPacketForUnknownThing(t *testing.T) {

    ctx := piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
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
}
