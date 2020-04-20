package main_test

import (
    "context"
    "fmt"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
)

func TestUsersGet(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    CreateUser(t, db, "user1@test.com", "")

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
    db := GetDb(t)
    CleanDb(t, db)
    orgId := CreateOrg(t, db, "org1")
    org2Id := CreateOrg(t, db, "org2")
    userId := CreateUser(t, db, "user1@test.com", "")
    AddOrgUser(t, db, orgId, userId)
    AddOrgUser(t, db, org2Id, userId)

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
    db := GetDb(t)
    CleanDb(t, db)

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
        Query: `
            mutation {
                createUser(email: "user_new@test.com", password: "pwd") { email }
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
    db := GetDb(t)
    CleanDb(t, db)
    id := CreateUser(t, db, "user1@test.com", "")

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
