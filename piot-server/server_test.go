package main

import (
    "fmt"
    "testing"
    "os"
    "encoding/json"
    "strings"
    "time"
    "context"
    "net/http"
    "net/http/httptest"
    "piot-server/handler"
    "piot-server/service"
    "piot-server/model"
    "piot-server/resolver"
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


/*
t.Run("gql users are protected", testGqlUsersNoAuthFunc(&ctx))
t.Run("gql users", testGqlUsersFunc(&ctx))
t.Run("gql create user", testGqlUserCreateFunc(&ctx, "test2@test.com"))
t.Run("gql update user", testGqlUserUpdateFunc(&ctx, "test2@test.com"))

t.Run("gql customers are protected", testGqlCustomersNoAuthFunc(&ctx))
t.Run("gql customers", testGqlCustomersFunc(&ctx))
*/

func login(t *testing.T, ctx *context.Context, email string, password string, statusCode int) (string) {
    req, err := http.NewRequest("POST", "/login", strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)))
    test.Ok(t, err)

    rr := httptest.NewRecorder()

    handler := handler.AddContext(*ctx, handler.LoginHandler())
    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, statusCode)

    var response model.Token
    test.Ok(t, json.Unmarshal(test.Body2Bytes(rr.Body), &response))

    return response.Token
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

func getAuthGqlRequest(t *testing.T, ctx *context.Context, body string) (*http.Request) {
    token := login(t, ctx, ADMIN_EMAIL, ADMIN_PASSWORD, 200)

    req, err := http.NewRequest("POST", "/any-path", strings.NewReader(body))
    test.Ok(t, err)
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

    return req
}

func TestLoginSuccessful(t *testing.T) {
    login(t, &ctx, ADMIN_EMAIL, ADMIN_PASSWORD, http.StatusOK)
}

func TestLoginWrongPassword(t *testing.T) {
    login(t, &ctx, ADMIN_EMAIL, "xxx", 401)
}

func TestLoginWrongEmail(t *testing.T) {
    login(t, &ctx, "xxx", ADMIN_PASSWORD, 401)
}

func TestLoginWrongEmailAndPassword(t *testing.T) {
    login(t, &ctx, "xxx", "yyy", 401)
}

func TestLoginEmptyEmailAndPassword(t *testing.T) {
    login(t, &ctx, "", "", 401)
}

func TestGqlUsersNoAuth(t *testing.T) {
    req, err := http.NewRequest("POST", "/any-path", strings.NewReader(`{"query":"{users {email}}"}`))
    if err != nil {t.Fatal(err)}

    rr := httptest.NewRecorder()

    // create GraphQL schema
    graphqlSchema := graphql.MustParseSchema(GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}

func TestGqlUsersGet(t *testing.T) {
    req := getAuthGqlRequest(t, &ctx, `{"query":"{users {email}}"}`)

    rr := httptest.NewRecorder()

    // create GraphQL schema
    graphqlSchema := graphql.MustParseSchema(GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 200)
}

func TestGqlUserCreate(t *testing.T) {

    const email = "test2@test.com"

    req := getAuthGqlRequest(t, &ctx, fmt.Sprintf(`{"query":"mutation {createUser(email: \"%s\") {id} }"}`, email))

    rr := httptest.NewRecorder()

    graphqlSchema := graphql.MustParseSchema(GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    test.CheckGqlResult(t, rr)
}

func TestGqlUserUpdate(t *testing.T) {

    const email = "test_create@test.com"
    const emailNew = "test_create_new@test.com"

    // create user
    id := createUser(t, &ctx, email, "pwd")
    t.Logf("User to be updated %s", id)

    // update user created in prev. step
    req := getAuthGqlRequest(t, &ctx, fmt.Sprintf(`{"query":"mutation {updateUser(id: \"%s\", email: \"%s\") {id} }"}`, id, emailNew))

    rr := httptest.NewRecorder()

    graphqlSchema := graphql.MustParseSchema(GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    test.CheckGqlResult(t, rr)

    // try to get user based on updated email address
    getUser(t, &ctx, emailNew)
}

func TestGqlCustomersNoAuth(t *testing.T) {
    req, err := http.NewRequest("POST", "/any-path", strings.NewReader(`{"query":"{customers {id}}"}`))
    if err != nil {t.Fatal(err)}

    rr := httptest.NewRecorder()

    graphqlSchema := graphql.MustParseSchema(GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 401)
}

func TestGqlCustomersGet(t *testing.T) {

    req := getAuthGqlRequest(t, &ctx, `{"query":"{customers{id}}"}`)

    rr := httptest.NewRecorder()

    // create GraphQL schema
    graphqlSchema := graphql.MustParseSchema(GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    test.CheckStatusCode(t, rr, 200)
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
