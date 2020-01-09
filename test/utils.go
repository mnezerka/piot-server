package test

import (
    "context"
    "fmt"
    "path/filepath"
    "net/http/httptest"
    "os"
    "runtime"
    "reflect"
    "testing"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "piot-server/utils"
    piotcontext "piot-server/context"
)

func CreateTestContext() context.Context {
    contextOptions := piotcontext.NewContextOptions()
    contextOptions.DbUri = os.Getenv("MONGODB_URI")
    contextOptions.DbName = "piot-test"
    contextOptions.MqttUri = "mock"
    contextOptions.Params.LogLevel = "DEBUG"
    ctx := piotcontext.NewContext(contextOptions)

    callerEmail := "caller@test.com"
    ctx = context.WithValue(ctx, "user_email", &callerEmail)
    ctx = context.WithValue(ctx, "is_authorized", true)
    return ctx
}

// assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
    if !condition {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
        tb.FailNow()
    }
}

// ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
    if err != nil {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
        tb.FailNow()
    }
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
    if !reflect.DeepEqual(exp, act) {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d:\n\texp: %#v\n\tgot: %#v\033[39m\n", filepath.Base(file), line, exp, act)
        tb.FailNow()
    }
}

// helper function for checking and logging respone status
func CheckStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
    if status := rr.Code; status != expected {
        t.Errorf("\033[31mWrong response status code: got %v want %v, body:\n%s\033[39m",
            status, expected, rr.Body.String())
    }
}

func CleanDb(t *testing.T, ctx context.Context) {
    db := ctx.Value("db").(*mongo.Database)
    db.Collection("orgs").DeleteMany(ctx, bson.M{})
    db.Collection("users").DeleteMany(ctx, bson.M{})
    db.Collection("orgusers").DeleteMany(ctx, bson.M{})
    db.Collection("things").DeleteMany(ctx, bson.M{})
    t.Log("DB is clean")
}

func CreateThing(t *testing.T, ctx context.Context, name string) (primitive.ObjectID) {

    db := ctx.Value("db").(*mongo.Database)

    res, err := db.Collection("things").InsertOne(ctx, bson.M{
        "name": name,
        "type": "sensor",
        "created": int32(time.Now().Unix()),
        "enabled": true,
        "sensor": bson.M{
            "class": "temperature",
            "measurement_topic": "value",
        },
    })
    Ok(t, err)

    t.Logf("Created thing %v", res.InsertedID)

    return res.InsertedID.(primitive.ObjectID)
}

func CreateUser(t *testing.T, ctx context.Context, email, password string) (primitive.ObjectID) {

    db := ctx.Value("db").(*mongo.Database)

    hash, err := utils.GetPasswordHash(password)
    Ok(t, err)

    res, err := db.Collection("users").InsertOne(ctx, bson.M{
        "email": email,
        "password": hash,
        "created": int32(time.Now().Unix()),
    })
    Ok(t, err)

    t.Logf("Created user %v", res.InsertedID)

    return res.InsertedID.(primitive.ObjectID)
}

func CreateOrg(t *testing.T, ctx context.Context, name string) (primitive.ObjectID) {
    db := ctx.Value("db").(*mongo.Database)

    res, err := db.Collection("orgs").InsertOne(ctx, bson.M{
        "name": name,
        "created": int32(time.Now().Unix()),
        "influxdb": "db",
        "influxdb_username": "db-username",
        "influxdb_password": "db-password",
    })
    Ok(t, err)

    t.Logf("Created org %v", res.InsertedID)

    return res.InsertedID.(primitive.ObjectID)
}

func AddOrgUser(t *testing.T, ctx context.Context, orgId, userId primitive.ObjectID) {
    db := ctx.Value("db").(*mongo.Database)

    _, err := db.Collection("orgusers").InsertOne(ctx, bson.M{
        "org_id": orgId,
        "user_id": userId,
        "created": int32(time.Now().Unix()),
    })
    Ok(t, err)

    t.Logf("User %v added to org %v", userId.Hex(), orgId.Hex())
}

func AddOrgThing(t *testing.T, ctx context.Context, orgId primitive.ObjectID, thingName string) {
    db := ctx.Value("db").(*mongo.Database)

    _, err := db.Collection("things").UpdateOne(ctx, bson.M{"name": thingName}, bson.M{"$set": bson.M{"org_id": orgId}})
    Ok(t, err)

    t.Logf("Thing %s assigned to org %s", thingName, orgId.Hex())
}

