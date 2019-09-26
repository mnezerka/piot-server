package resolver

import (
    "context"
    "fmt"
    "testing"
    "time"
    "os"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
    "piot-server/test"
    piotcontext "piot-server/context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"

)

var customerId string

func createCustomer(t *testing.T, ctx *context.Context, name string) (string) {

    db := (*ctx).Value("db").(*mongo.Database)

    res, err := db.Collection("customers").InsertOne(*ctx, bson.M{
        "name": name,
        "created": int32(time.Now().Unix()),
    })
    test.Ok(t, err)

    return res.InsertedID.(primitive.ObjectID).Hex()
}


func cleanDb(t *testing.T, ctx context.Context) {

    db := ctx.Value("db").(*mongo.Database)
    db.Collection("customers").DeleteMany(ctx, bson.M{})
    db.Collection("users").DeleteMany(ctx, bson.M{})
}

func init() {

    ctx = piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    ctx = context.WithValue(ctx, "user_email", "admin@test.com")
    ctx = context.WithValue(ctx, "is_authorized", true)

    // close database
    //ctx.Value("dbClient").(*mongo.Client).Disconnect(ctx)
}

func TestCustomers(t *testing.T) {
    cleanDb(t, ctx)
    createCustomer(t, &ctx, "customer1")

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: ctx,
            Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
            Query: `
                {
                    customers { name }
                }
            `,
            ExpectedResult: `
                {
                    "customers": [
                        {
                            "name": "customer1"
                        }
                    ]
                }
            `,
        },
    })
}

func TestCustomer(t *testing.T) {
    cleanDb(t, ctx)
    customerId = createCustomer(t, &ctx, "customer1")

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            {
                customer(id: "%s") { name }
            }
        `, customerId),
        ExpectedResult: `
            {
                "customer": {
                    "name": "customer1"
                }
            }
        `,
    })
}
