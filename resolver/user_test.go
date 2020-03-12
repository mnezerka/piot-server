package resolver_test

import (
    "context"
    "fmt"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "github.com/mnezerka/go-piot/test"
    "piot-server/schema"
)

func TestUsersGet(t *testing.T) {
    db := test.GetDb(t)
    test.CleanDb(t, db)
    test.CreateUser(t, db, "user1@test.com", "")

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: context.TODO(),
            Schema: schema,
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
    db := test.GetDb(t)
    test.CleanDb(t, db)
    orgId := test.CreateOrg(t, db, "org1")
    org2Id := test.CreateOrg(t, db, "org2")
    userId := test.CreateUser(t, db, "user1@test.com", "")
    test.AddOrgUser(t, db, orgId, userId)
    test.AddOrgUser(t, db, org2Id, userId)

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema:  schema,
        Query: fmt.Sprintf(`
            {
                user(id: "%s") {email, orgs {name}}
            }
        `, userId.Hex()),
        ExpectedResult: `
            {
                "user": {
                    "email": "user1@test.com",
                    "orgs": [{"name": "org1"}, {"name": "org2"}]
                }
            }
        `,
    })
}

func TestUserCreate(t *testing.T) {
    db := test.GetDb(t)
    test.CleanDb(t, db)

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
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
    db := test.GetDb(t)
    test.CleanDb(t, db)
    id := test.CreateUser(t, db, "user1@test.com", "")

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    t.Logf("User to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
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
