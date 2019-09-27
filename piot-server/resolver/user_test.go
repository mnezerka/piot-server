package resolver

import (
    "context"
    "fmt"
    "os"
    "testing"
    "time"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
    "piot-server/test"
    piotcontext "piot-server/context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx context.Context

func createUser(t *testing.T, ctx *context.Context, email string) (string) {

    db := (*ctx).Value("db").(*mongo.Database)

    res, err := db.Collection("users").InsertOne(*ctx, bson.M{
        "email": email,
        "created": int32(time.Now().Unix()),
    })
    test.Ok(t, err)

    return res.InsertedID.(primitive.ObjectID).Hex()
}

func init() {
    callerEmail := "caller@test.com"
    ctx = piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    ctx = context.WithValue(ctx, "user_email", &callerEmail)
    ctx = context.WithValue(ctx, "is_authorized", true)
}

func TestUsersGet(t *testing.T) {
    cleanDb(t, ctx)
    createUser(t, &ctx, "user1@test.com")

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: ctx,
            Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
            Query: `
                {
                    users { email }
                }
            `,
            ExpectedResult: `
                {
                    "users": [
                        {
                            "email": "user1@test.com"
                        }
                    ]
                }
            `,
        },
    })
}

func TestUserGet(t *testing.T) {
    cleanDb(t, ctx)
    id := createUser(t, &ctx, "user1@test.com")

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            {
                user(id: "%s") { email, org { id } }
            }
        `, id),
        ExpectedResult: `
            {
                "user": {
                    "email": "user1@test.com",
                    "org": null
                }
            }
        `,
    })
}

func TestUserCreate(t *testing.T) {
    cleanDb(t, ctx)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: `
            mutation {
                createUser(user: {email: "user_new@test.com"}) { email }
            }
        `,
        ExpectedResult: `
            {
                "createUser": {
                    "email": "user_new@test.com"
                }
            }
        `,
    })
}

func TestUserUpdate(t *testing.T) {
    cleanDb(t, ctx)
    id := createUser(t, &ctx, "user1@test.com")

    t.Logf("User to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            mutation {
                updateUser(user: {id: "%s", email: "user1_new@test.com"}) { email }
            }
        `, id),
        ExpectedResult: `
            {
                "updateUser": {
                    "email": "user1_new@test.com"
                }
            }
        `,
    })
}
