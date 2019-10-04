package handler_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "piot-server/test"
    "piot-server/handler"
)

func TestAuthUser(t *testing.T) {
    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

    deviceData := `
    {
        "username": "xxx",
        "password": "xxx"
    }`

    req, err := http.NewRequest("POST", "/mosquitto-auth-user", strings.NewReader(deviceData))
    test.Ok(t, err)

    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()

    handler := handler.MosquittoAuth{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}



func TestAuthAcl(t *testing.T) {
    ctx := test.CreateTestContext()

    test.CleanDb(t, ctx)

    deviceData := `
    {
        "acc":2,
        "clientid":"mqtt-client",
        "topic":"hello",
        "username":"test@test.com"
    }`

    req, err := http.NewRequest("POST", "/mosquitto-auth-acl", strings.NewReader(deviceData))
    test.Ok(t, err)

    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()

    handler := handler.MosquittoAuth{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}

func TestAuthSuperUser(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)

    deviceData := `
    {
        "username": "xxx",
    }`

    req, err := http.NewRequest("POST", "/mosquitto-auth-superuser", strings.NewReader(deviceData))
    test.Ok(t, err)

    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()

    handler := handler.MosquittoAuth{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}

func TestAuthInvalidPath(t *testing.T) {
    ctx := test.CreateTestContext()

    deviceData := "something as body"

    req, err := http.NewRequest("POST", "/mosquitto-auth-user-xx", strings.NewReader(deviceData))
    test.Ok(t, err)

    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()

    handler := handler.MosquittoAuth{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 403)
}

