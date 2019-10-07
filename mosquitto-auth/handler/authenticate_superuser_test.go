package handler_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "mosquitto-auth/test"
    "mosquitto-auth/handler"
)

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
