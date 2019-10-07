package handler_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "mosquitto-auth/test"
    "mosquitto-auth/handler"
)

func TestAuthUser(t *testing.T) {
    ctx := test.CreateTestContext()

    deviceData := `
    {
        "username": "xxx",
        "password": "xxx"
    }`

    req, err := http.NewRequest("POST", "/mosquitto-auth-user", strings.NewReader(deviceData))
    test.Ok(t, err)

    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()

    handler := handler.AuthenticateUser{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}

func TestAuthAcl(t *testing.T) {
    ctx := test.CreateTestContext()

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

    handler := handler.Authorize{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}

func TestAuthSuperUser(t *testing.T) {
    ctx := test.CreateTestContext()

    deviceData := `
    {
        "username": "xxx",
    }`

    req, err := http.NewRequest("POST", "/mosquitto-auth-superuser", strings.NewReader(deviceData))
    test.Ok(t, err)

    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()

    handler := handler.AuthenticateSuperUser{}

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}
