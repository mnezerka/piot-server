package resolver

/*
import (
    "context"
    "os"
    "testing"
    "time"
    piotcontext "piot-server/context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)
*/
/*
var ctx context.Context

var customerId string

func createCustomer(ctx *context.Context, name string) (string) {

    db := (*ctx).Value("db").(*mongo.Database)

    res, err := db.Collection("customers").InsertOne(*ctx, bson.M{
        "name": name,
        "created": int32(time.Now().Unix()),
    })
    if err != nil {
        os.Exit(1)
    }

    return res.InsertedID.(primitive.ObjectID).Hex()
}

func TestMain(m *testing.M) {

    // prepare context
    ctx = piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    db := ctx.Value("db").(*mongo.Database)
    ctx = context.WithValue(ctx, "user_email", "admin@test.com")
    ctx = context.WithValue(ctx, "is_authorized", true)

    // clean db
    db.Collection("customers").DeleteMany(ctx, bson.M{})
    db.Collection("users").DeleteMany(ctx, bson.M{})

    // create seed data
    //userId = createUser(&ctx, "user1@test.com")
    //customerId = createCustomer(&ctx, "customer1")

    os.Exit(m.Run())

    // close database
    ctx.Value("dbClient").(*mongo.Client).Disconnect(ctx)
}
*/
