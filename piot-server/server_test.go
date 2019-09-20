package main

import (
    "bytes"
    "testing"
    "os"
    "encoding/json"
    "io/ioutil"
    "strings"
    "time"
    "context"
    "net/http"
    "net/http/httptest"
    "piot-server/handler"
    "piot-server/service"
    "piot-server/model"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

func TestAPI(t *testing.T) {
    ctx := context.Background()

    // create global logger for all handlers
    log := service.NewLogger(LOG_FORMAT, true)
    ctx = context.WithValue(ctx, "log", log)

    // try to open database
    dbUri := os.Getenv("MONGODB_URI")
    dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
    FatalOnError(err, "Failed to open database on %s", dbUri)

    // Check the connection
    err = dbClient.Ping(ctx, nil)
    FatalOnError(err, "Cannot ping database on %s", dbUri)

    // Auto disconnect from mongo
    defer dbClient.Disconnect(ctx)

    db := dbClient.Database("piot-test")
    ctx = context.WithValue(ctx, "db", db)

    // create admin account
    hash, err := handler.GetPasswordHash("test")
    FatalOnError(err, "Cannot generate hash from password")

    db.Collection("users").DeleteMany(ctx, bson.M{})
    _, err = db.Collection("users").InsertOne(ctx, bson.M{
        "email": "test@test.com",
        "password": hash,
        "created": int32(time.Now().Unix()),
    })
    FatalOnError(err, "Cannot insert test user account")

    //////////// run tests

    t.Run("login", testLoginFunc(&ctx))
}

// helper function for checking and logging respone status
func checkStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
    if status := rr.Code; status != expected{

        t.Errorf("handler returned wrong status code: got %v want %v, body:\n%s",
            status, expected, rr.Body.String())
    }
}

func body2Bytes(body *bytes.Buffer) ([]byte) {

    var result []byte
    result, _ = ioutil.ReadAll(body)
    //result = byte[](res.Body.String())
    return result
}

func testLoginFunc(ctx *context.Context) func(*testing.T) {
    return func(t *testing.T) {
        req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"email": "test@test.com", "password": "test"}`))
        if err != nil {t.Fatal(err)}

        rr := httptest.NewRecorder()

        handler := handler.AddContext(*ctx, handler.LoginHandler())
        handler.ServeHTTP(rr, req)

        checkStatusCode(t, rr, http.StatusOK)

        var response model.Token
        err = json.Unmarshal(body2Bytes(rr.Body), &response)
        if err != nil {t.Fatal(err)}
    }
}

func TestRoot(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handler.RootHandler)
    handler.ServeHTTP(rr, req)
    checkStatusCode(t, rr, http.StatusOK)

    // Check the response body is what we expect.
    if !strings.HasPrefix(rr.Body.String(), "<html>") {
        t.Error("unexpected body: does start with <html>")
    }
}
