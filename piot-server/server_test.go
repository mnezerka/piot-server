package main

import (
    "bytes"
    "fmt"
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
    "piot-server/resolver"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    graphql "github.com/graph-gophers/graphql-go"
)

const ADMIN_EMAIL = "test@test.com"
const ADMIN_PASSWORD = "test"

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
    hash, err := handler.GetPasswordHash(ADMIN_PASSWORD)
    FatalOnError(err, "Cannot generate hash from password")

    db.Collection("users").DeleteMany(ctx, bson.M{})
    _, err = db.Collection("users").InsertOne(ctx, bson.M{
        "email": ADMIN_EMAIL,
        "password": hash,
        "created": int32(time.Now().Unix()),
    })
    FatalOnError(err, "Cannot insert test user account")

    //////////// run tests

    t.Run("login ok", testLoginFunc(&ctx, ADMIN_EMAIL, ADMIN_PASSWORD, http.StatusOK))
    t.Run("login wrong password", testLoginFunc(&ctx, ADMIN_EMAIL, "xxx", 401))
    t.Run("login wrong email", testLoginFunc(&ctx, "xxx", ADMIN_PASSWORD, 401))
    t.Run("login wrong email and password", testLoginFunc(&ctx, "xxx", "yyy", 401))
    t.Run("login empty email and password", testLoginFunc(&ctx, "", "", 401))

    t.Run("gql users are protected", testGqlUsersFunc(&ctx))
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
    return result
}

func testLoginFunc(ctx *context.Context, email string, password string, statusCode int) func(*testing.T) {
    return func(t *testing.T) {
        req, err := http.NewRequest("POST", "/login", strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)))
        if err != nil {t.Fatal(err)}

        rr := httptest.NewRecorder()

        handler := handler.AddContext(*ctx, handler.LoginHandler())
        handler.ServeHTTP(rr, req)

        checkStatusCode(t, rr, statusCode)

        var response model.Token
        err = json.Unmarshal(body2Bytes(rr.Body), &response)
        if err != nil {t.Fatal(err)}
    }
}

func testGqlUsersFunc(ctx *context.Context) func(*testing.T) {
    return func(t *testing.T) {
        req, err := http.NewRequest("POST", "/any-path", strings.NewReader(`{"query":"{users {email}}"}`))
        if err != nil {t.Fatal(err)}

        rr := httptest.NewRecorder()

        // create GraphQL schema
        graphqlSchema := graphql.MustParseSchema(GetRootSchema(), &resolver.Resolver{})

        handler := handler.AddContext(*ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

        handler.ServeHTTP(rr, req)

        checkStatusCode(t, rr, 401)
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
