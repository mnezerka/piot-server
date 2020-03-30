package main_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "piot-server"
    "piot-server/model"
)

type mockHandlerCall struct {
    Request *http.Request
}

type mockHandler struct {
    Log *logging.Logger
    Calls []mockHandlerCall
}

func getMockHandler(logger *logging.Logger) *mockHandler {
    h := &mockHandler{}
    h.Log = logger
    return h
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.Log.Debugf("Mock handler called")
    h.Calls = append(h.Calls, mockHandlerCall{r})
}

func getAuthHandler(t *testing.T, db *mongo.Database, h http.Handler) *main.AuthHandler {
    log := GetLogger(t)
    cfg := GetConfig()
    users := GetUsers(t, log, db)
    return main.NewAuthHandler(log, cfg, users, h)
}

// Missing and invalid authorization data
func TestAuthNoCredentials(t *testing.T) {
    log := GetLogger(t)
    db := GetDb(t)

    // request without headers
    req, err := http.NewRequest("POST", "/", strings.NewReader(""))
    Ok(t, err)

    rr := httptest.NewRecorder()

    handler := getAuthHandler(t, db, getMockHandler(log))
    handler.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 401)


    // request with invalid authorization header
    req, err = http.NewRequest("POST", "/", strings.NewReader(""))
    req.Header.Add("Auhthorization", "XXX")
    Ok(t, err)

    rr = httptest.NewRecorder()

    handler = getAuthHandler(t, db, getMockHandler(log))
    handler.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 401)
}

// Authenticated valid user
func TestAuthValid(t *testing.T) {
    log := GetLogger(t)
    db := GetDb(t)
    CleanDb(t, db)
    userId := CreateUser(t, db, ADMIN_EMAIL, ADMIN_PASSWORD)
    orgId := CreateOrg(t, db, "Org")
    AddOrgUser(t, db, orgId, userId)

    token := LoginUser(t, log, db, ADMIN_EMAIL, ADMIN_PASSWORD, http.StatusOK)

    // send some request and let handler to initiate user profile section of
    // context associated with request
    req, err := http.NewRequest("POST", "/", strings.NewReader(""))
    req.Header.Add("Authorization", "Bearer " + token)

    Ok(t, err)

    rr := httptest.NewRecorder()

    mh := getMockHandler(log)

    handler := getAuthHandler(t, db, mh)
    handler.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 200)

    // check if child handler was called
    Equals(t, 1, len(mh.Calls))

    // get context associated with child request
    ctx := mh.Calls[0].Request.Context()

    // verify that context contains user profile
    profile := ctx.Value("profile").(*model.UserProfile)

    Assert(t, profile != nil, "User profile not initialized")
    Equals(t, profile.Email, ADMIN_EMAIL)
    Equals(t, profile.IsAdmin, false)
    //Assert(t, profile.Org != nil, "User profile has no org")
    Equals(t, 1, len(profile.Orgs))
}
