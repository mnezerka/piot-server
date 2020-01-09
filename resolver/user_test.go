package resolver

import (
    "fmt"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
    "piot-server/test"
)

func TestUsersGet(t *testing.T) {
    ctx := test.CreateTestContext()
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
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    orgId := test.CreateOrg(t, ctx, "org1")
    userId := test.CreateUser(t, ctx, "user1@test.com", "")
    test.AddOrgUser(t, ctx, orgId, userId)

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
    ctx := test.CreateTestContext()
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
    ctx := test.CreateTestContext()
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
