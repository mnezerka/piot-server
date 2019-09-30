package handler_test

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http/httptest"
    "strings"
    "testing"
    "piot-server/model"
    "piot-server/test"
    "piot-server/handler"
    "net/http"
)

var ADMIN_EMAIL = "admin@test.com"
var ADMIN_PASSWORD = "admin"

const TEST_EMAIL = "test@test.com"
const TEST_PASSWORD = "test"

func Body2Bytes(body *bytes.Buffer) ([]byte) {
    var result []byte
    result, _ = ioutil.ReadAll(body)
    return result
}

func Login(t *testing.T, ctx *context.Context, email string, password string, statusCode int) (string) {
    req, err := http.NewRequest("POST", "/login", strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)))
    test.Ok(t, err)

    rr := httptest.NewRecorder()

    handler := handler.AddContext(*ctx, handler.LoginHandler())
    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, statusCode)

    var response model.Token
    test.Ok(t, json.Unmarshal(Body2Bytes(rr.Body), &response))

    return response.Token
}
