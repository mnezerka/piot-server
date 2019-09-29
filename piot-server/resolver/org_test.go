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

    return res.InsertedID.(primitive.ObjectID)
}

func AssignOrgUser(t *testing.T, ctx *context.Context, orgId, userId primitive.ObjectID) {
    db := (*ctx).Value("db").(*mongo.Database)

    /*
    orgId, err := primitive.ObjectIDFromHex(string(args.OrgId))
    test.Ok(t.err)

    userId, err := primitive.ObjectIDFromHex(string(args.OrgId))
    */

    _, err := db.Collection("orgusers").InsertOne(*ctx, bson.M{
        "org_id": orgId,
        "user_id": userId,
        "created": int32(time.Now().Unix()),
    })
    test.Ok(t, err)

    t.Logf("Assigned user %v to org %v", userId.Hex(), orgId.Hex())
}

func CleanDb(t *testing.T, ctx context.Context) {
    db := ctx.Value("db").(*mongo.Database)
    db.Collection("orgs").DeleteMany(ctx, bson.M{})
    db.Collection("users").DeleteMany(ctx, bson.M{})
    db.Collection("orgusers").DeleteMany(ctx, bson.M{})
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
    AssignOrgUser(t, &ctx, orgId, userId)

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

func TestAssignOrgUser(t *testing.T) {
    CleanDb(t, ctx)
    userId := CreateUser(t, &ctx, "user1@test.com")
    orgId := CreateOrg(t, &ctx, "test-org")

    t.Logf("User to be assigned %s, org to be assigned %s", userId, orgId)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            mutation {
                assignOrgUser(orgId: "%s", userId: "%s")
            }
        `, orgId, userId),
        ExpectedResult: `
            {
                "assignOrgUser": null
            }
        `,
    })
}

