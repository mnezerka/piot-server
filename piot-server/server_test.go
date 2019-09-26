package main

import (
    "fmt"
    "testing"
    "os"
//    "encoding/json"
    "strings"
    "time"
    "context"
    "net/http"
    "net/http/httptest"
    "piot-server/handler"
    "piot-server/service"
    "piot-server/model"
    "piot-server/resolver"
    "piot-server/schema"
    "piot-server/test"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    graphql "github.com/graph-gophers/graphql-go"
)

const ADMIN_EMAIL = "admin@test.com"
const ADMIN_PASSWORD = "admin"

const TEST_EMAIL = "test@test.com"
const TEST_PASSWORD = "test"

var ctx context.Context

func TestMain(m *testing.M) {

    ctx = context.Background()

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

    // Users
    db.Collection("users").DeleteMany(ctx, bson.M{})

    // create admin account
    hash, err := handler.GetPasswordHash(ADMIN_PASSWORD)
    FatalOnError(err, "Cannot generate hash from password")

    _, err = db.Collection("users").InsertOne(ctx, bson.M{
        "email": ADMIN_EMAIL,
        "password": hash,
        "created": int32(time.Now().Unix()),
    })
    FatalOnError(err, "Cannot insert admin user account")

    hash, err = handler.GetPasswordHash(TEST_PASSWORD)
    FatalOnError(err, "Cannot generate hash from password")

    _, err = db.Collection("users").InsertOne(ctx, bson.M{
        "email": TEST_EMAIL,
        "password": hash,
        "created": int32(time.Now().Unix()),
    })
    FatalOnError(err, "Cannot insert test user account")


    os.Exit(m.Run())
}

func createUser(t *testing.T, ctx *context.Context, email string, password string) (string) {

    db := (*ctx).Value("db").(*mongo.Database)

    hash, err := handler.GetPasswordHash(password)
    test.Ok(t, err)

    res, err := db.Collection("users").InsertOne(*ctx, bson.M{
        "email": email,
        "password": hash,
        "created": int32(time.Now().Unix()),
    })
    test.Ok(t, err)

    return res.InsertedID.(primitive.ObjectID).Hex()
}

func getUser(t *testing.T, ctx *context.Context, email string) (*model.User) {
    db := (*ctx).Value("db").(*mongo.Database)

    var user model.User
    err := db.Collection("users").FindOne(*ctx, bson.D{{"email", email}}).Decode(&user)
    test.Ok(t, err)

    return &user
}

func TestLoginSuccessful(t *testing.T) {
    test.Login(t, &ctx, ADMIN_EMAIL, ADMIN_PASSWORD, http.StatusOK)
}

func TestLoginWrongPassword(t *testing.T) {
    test.Login(t, &ctx, ADMIN_EMAIL, "xxx", 401)
}

func TestLoginWrongEmail(t *testing.T) {
    test.Login(t, &ctx, "xxx", ADMIN_PASSWORD, 401)
}

func TestLoginWrongEmailAndPassword(t *testing.T) {
    test.Login(t, &ctx, "xxx", "yyy", 401)
}

func TestLoginEmptyEmailAndPassword(t *testing.T) {
    test.Login(t, &ctx, "", "", 401)
}

func TestGqlUsersNoAuth(t *testing.T) {
    req, err := http.NewRequest("POST", "/any-path", strings.NewReader(`{"query":"{users {email}}"}`))
    test.Ok(t, err)

    rr := httptest.NewRecorder()

    graphqlSchema := graphql.MustParseSchema(schema.GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}

func TestGqlUsersGet(t *testing.T) {

    rr := test.GetGqlResponseRecorder(t, &ctx, ADMIN_EMAIL, ADMIN_PASSWORD, `{"query":"{users {email, customer{id}}}"}`)

    test.CheckGqlResult(t, rr)
}

func TestGqlUserCreate(t *testing.T) {

    const email = "test2@test.com"

    request := fmt.Sprintf(`{"query":"mutation {createUser(email: \"%s\") {id} }"}`, email)

    rr := test.GetGqlResponseRecorder(t, &ctx, ADMIN_EMAIL, ADMIN_PASSWORD, request)

    test.CheckGqlResult(t, rr)
}

func TestGqlUserUpdate(t *testing.T) {

    const email = "test_create@test.com"
    const emailNew = "test_create_new@test.com"

    // create user
    id := createUser(t, &ctx, email, "pwd")
    t.Logf("User to be updated %s", id)

    // update user created in prev. step
    request := fmt.Sprintf(`{"query":"mutation {updateUser(id: \"%s\", email: \"%s\") {id} }"}`, id, emailNew)

    rr := test.GetGqlResponseRecorder(t, &ctx, ADMIN_EMAIL, ADMIN_PASSWORD, request)

    test.CheckGqlResult(t, rr)

    // try to get user based on updated email address
    getUser(t, &ctx, emailNew)
}

func TestGqlCustomersNoAuth(t *testing.T) {
    req, err := http.NewRequest("POST", "/any-path", strings.NewReader(`{"query":"{customers {id}}"}`))
    test.Ok(t, err)

    rr := httptest.NewRecorder()

    graphqlSchema := graphql.MustParseSchema(schema.GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}

func TestGqlCustomersGet(t *testing.T) {

    rr := test.GetGqlResponseRecorder(t, &ctx, ADMIN_EMAIL, ADMIN_PASSWORD, `{"query":"{customers{id}}"}`)
    test.CheckGqlResult(t, rr)
}

func TestRoot(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handler.RootHandler)
    handler.ServeHTTP(rr, req)
    test.CheckStatusCode(t, rr, http.StatusOK)

    // Check the response body is what we expect.
    if !strings.HasPrefix(rr.Body.String(), "<html>") {
        t.Error("unexpected body: does start with <html>")
    }
}
