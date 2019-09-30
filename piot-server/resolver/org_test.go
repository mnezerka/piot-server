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

func CreateOrg(t *testing.T, ctx *context.Context, name string) (primitive.ObjectID) {
    db := (*ctx).Value("db").(*mongo.Database)

    res, err := db.Collection("orgs").InsertOne(*ctx, bson.M{
        "name": name,
        "created": int32(time.Now().Unix()),
    })
    test.Ok(t, err)

    t.Logf("Created org %v", res.InsertedID)

    return res.InsertedID.(primitive.ObjectID)
}

func AddOrgUser(t *testing.T, ctx *context.Context, orgId, userId primitive.ObjectID) {
    db := (*ctx).Value("db").(*mongo.Database)

    _, err := db.Collection("orgusers").InsertOne(*ctx, bson.M{
        "org_id": orgId,
        "user_id": userId,
        "created": int32(time.Now().Unix()),
    })
    test.Ok(t, err)

    t.Logf("User %v added to org %v", userId.Hex(), orgId.Hex())
}

func CleanDb(t *testing.T, ctx context.Context) {
    db := ctx.Value("db").(*mongo.Database)
    db.Collection("orgs").DeleteMany(ctx, bson.M{})
    db.Collection("users").DeleteMany(ctx, bson.M{})
    db.Collection("orgusers").DeleteMany(ctx, bson.M{})

    t.Log("DB is clean")
}

func init() {
    ctx = piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    ctx = context.WithValue(ctx, "user_email", "admin@test.com")
    ctx = context.WithValue(ctx, "is_authorized", true)

    // close database
    //ctx.Value("dbClient").(*mongo.Client).Disconnect(ctx)
}

func TestOrgsGet(t *testing.T) {
    CleanDb(t, ctx)
    CreateOrg(t, &ctx, "org1")

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: ctx,
            Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
            Query: `
                {
                    orgs { name }
                }
            `,
            ExpectedResult: `
                {
                    "orgs": [
                        {
                            "name": "org1"
                        }
                    ]
                }
            `,
        },
    })
}

func TestOrgGet(t *testing.T) {
    CleanDb(t, ctx)
    orgId := CreateOrg(t, &ctx, "org1")
    userId := CreateUser(t, &ctx, "org1user@test.com")
    AddOrgUser(t, &ctx, orgId, userId)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            {
                org(id: "%s") {name, users {email}}
            }
        `, orgId.Hex()),
        ExpectedResult: `
            {
                "org": {
                    "name": "org1",
                    "users": [{"email": "org1user@test.com"}]
                }
            }
        `,
    })
}

func TestAddOrgUser(t *testing.T) {
    CleanDb(t, ctx)
    userId := CreateUser(t, &ctx, "user1@test.com")
    orgId := CreateOrg(t, &ctx, "test-org")
    org2Id := CreateOrg(t, &ctx, "test-org2")
    CreateOrg(t, &ctx, "test-org3")

    t.Logf("Test adding user %s to org %s", userId, orgId)

    // assign user to the first organization
    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            mutation {
                addOrgUser(orgId: "%s", userId: "%s")
            }
        `, orgId.Hex(), userId.Hex()),
        ExpectedResult: `
            {
                "addOrgUser": null
            }
        `,
    })

    t.Logf("Test adding user %s to org %s", userId, org2Id)

    // assign user to the second organization
    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            mutation {
                addOrgUser(orgId: "%s", userId: "%s")
            }
        `, org2Id.Hex(), userId.Hex()),
        ExpectedResult: `
            {
                "addOrgUser": null
            }
        `,
    })
}

func TestRemoveOrgUser(t *testing.T) {
    CleanDb(t, ctx)
    userId := CreateUser(t, &ctx, "user1@test.com")
    orgId := CreateOrg(t, &ctx, "test-org")
    org2Id := CreateOrg(t, &ctx, "test-org2")
    AddOrgUser(t, &ctx, orgId, userId)
    AddOrgUser(t, &ctx, org2Id, userId)

    t.Logf("Test remove user %s from org %s", userId, orgId)

    // assign user to the first organization
    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            mutation {
                removeOrgUser(orgId: "%s", userId: "%s")
            }
        `, orgId.Hex(), userId.Hex()),
        ExpectedResult: `
            {
                "removeOrgUser": null
            }
        `,
    })

    // TODO: check if user is still  
}
