package resolver

import (
    "context"
    "fmt"
    "os"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
    "piot-server/test"
    piotcontext "piot-server/context"
)

var ctx context.Context


func init() {
    callerEmail := "caller@test.com"
    ctx = piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    ctx = context.WithValue(ctx, "user_email", &callerEmail)
    ctx = context.WithValue(ctx, "is_authorized", true)
}

func TestUsersGet(t *testing.T) {
    test.CleanDb(t, ctx)
    test.CreateUser(t, ctx, "user1@test.com", "")

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
    ctx := piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    test.CleanDb(t, ctx)
    orgId := CreateOrg(t, &ctx, "org1")
    userId := test.CreateUser(t, ctx, "user1@test.com", "")
    AddOrgUser(t, &ctx, orgId, userId)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            {
                user(id: "%s") {email, orgs {name}}
            }
        `, userId.Hex()),
        ExpectedResult: `
            {
                "user": {
                    "email": "user1@test.com",
                    "orgs": [{"name": "org1"}]
                }
            }
        `,
    })
}

func TestUserCreate(t *testing.T) {
    test.CleanDb(t, ctx)

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
    test.CleanDb(t, ctx)
    id := test.CreateUser(t, ctx, "user1@test.com", "")

    t.Logf("User to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            mutation {
                updateUser(user: {id: "%s", email: "user1_new@test.com"}) { email }
            }
        `, id.Hex()),
        ExpectedResult: `
            {
                "updateUser": {
                    "email": "user1_new@test.com"
                }
            }
        `,
    })
}
