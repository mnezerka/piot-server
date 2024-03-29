package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	main "piot-server"
	"strings"
	"testing"

	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/mongo"
)

var ADMIN_EMAIL = "admin@test.com"
var ADMIN_PASSWORD = "admin"

const TEST_EMAIL = "test@test.com"
const TEST_PASSWORD = "test"

func Body2Bytes(body *bytes.Buffer) []byte {
	var result []byte
	result, _ = ioutil.ReadAll(body)
	return result
}

// helper function for checking and logging respone status
func CheckStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	if status := rr.Code; status != expected {
		t.Errorf("\033[31mWrong response status code: got %v want %v, body:\n%s\033[39m",
			status, expected, rr.Body.String())
	}
}

func LoginUser(t *testing.T, log *logging.Logger, db *mongo.Database, email string, password string, statusCode int) string {
	req, err := http.NewRequest("POST", "/login", strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)))
	Ok(t, err)

	rr := httptest.NewRecorder()

	//handler := handler.AddContext(*ctx, handler.LoginHandler())
	handler := main.NewLoginHandler(log, db, GetConfig())

	handler.ServeHTTP(rr, req)

	CheckStatusCode(t, rr, statusCode)

	var response main.Token
	Ok(t, json.Unmarshal(Body2Bytes(rr.Body), &response))

	return response.Token
}
