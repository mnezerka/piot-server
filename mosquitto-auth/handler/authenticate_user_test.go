package handler_test

import (
    "context"
    "fmt"
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

func TestAuthenticateStaticUsers(t *testing.T) {
    ctx := test.CreateTestContext()

    for _, username := range []string{"test", "mon", "piot"} {

        deviceData := fmt.Sprintf(`
        {
            "username": "%s",
            "password": "pwd-%s"
        }`, username, username)

        req, err := http.NewRequest("POST", "/", strings.NewReader(deviceData))
        test.Ok(t, err)

        // first request is with inactive static test user (empty password)
        req = req.WithContext(ctx)
        rr := httptest.NewRecorder()
        handler := handler.AuthenticateUser{}
        handler.ServeHTTP(rr, req)
        test.CheckStatusCode(t, rr, 401)

        // second request is with active static test user
        req, err = http.NewRequest("POST", "/", strings.NewReader(deviceData))
        test.Ok(t, err)

        ctx = context.WithValue(ctx, fmt.Sprintf("%s-pwd", username), fmt.Sprintf("pwd-%s", username))
        req = req.WithContext(ctx)
        rr = httptest.NewRecorder()
        handler.ServeHTTP(rr, req)
        test.CheckStatusCode(t, rr, 200)
    }
}
