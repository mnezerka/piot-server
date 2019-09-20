package main

import (
    "testing"
    "os"
    "strings"
    "context"
    "net/http"
    "net/http/httptest"
    "piot-server/handler"
    "piot-server/service"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
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

    db := dbClient.Database("piot")
    ctx = context.WithValue(ctx, "db", db)


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

func testLoginFunc(ctx *context.Context) func(*testing.T) {
    return func(t *testing.T) {
        req, err := http.NewRequest("POST", "/login", strings.NewReader("{}"))
        if err != nil { t.Fatal(err) }

        rr := httptest.NewRecorder()

        handler := handler.AddContext(*ctx, handler.LoginHandler())
        handler.ServeHTTP(rr, req)

        checkStatusCode(t, rr, http.StatusOK)

        // Check the response body is what we expect.
        if !strings.HasPrefix(rr.Body.String(), "<html>") {
            t.Error("unexpected body: does start with <html>")
        }
    }
}

func TestRoot(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    //handler := http.HandlerFunc(GetEntries)
    handler := http.HandlerFunc(handler.RootHandler)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Check the response body is what we expect.
    if !strings.HasPrefix(rr.Body.String(), "<html>") {
        t.Error("unexpected body: does start with <html>")
    }
}
