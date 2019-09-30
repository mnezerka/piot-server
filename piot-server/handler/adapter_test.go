package handler

import (
    //"context"
    "net/http"
    "net/http/httptest"
    //"fmt"
    "os"
    "strings"
    "testing"
    //"time"
    "piot-server/test"
    //"piot-server/model"
    piotcontext "piot-server/context"
    //"go.mongodb.org/mongo-driver/mongo"
    //"go.mongodb.org/mongo-driver/bson"
    //"go.mongodb.org/mongo-driver/bson/primitive"
)

//var ctx context.Context

func init() {
    //callerEmail := "caller@test.com"
    //ctx = piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    //ctx = context.WithValue(ctx, "user_email", &callerEmail)
    //ctx = context.WithValue(ctx, "is_authorized", true)
}

func TestForbiddenGet(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    test.Ok(t, err)
    req = req.WithContext(piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test"))

    rr := httptest.NewRecorder()

    handler := Adapter{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 405)
}

func TestUnknownDevice(t *testing.T) {

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
    req = req.WithContext(piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test"))

    rr := httptest.NewRecorder()

    handler := Adapter{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 200)
}
