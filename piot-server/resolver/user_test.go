package resolver

import (
    "context"
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

func TestUsers(t *testing.T) {
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

func TestUser(t *testing.T) {
    cleanDb(t, ctx)
    createUser(t, &ctx, "user1@test.com")

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: `
            {
                user(email: "user1@test.com") { email, customer { id } }
            }
        `,
        ExpectedResult: `
            {
                "user": {
                    "email": "user1@test.com",
                    "customer": null
                }
            }
        `,
    })
}


/*
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
*/

