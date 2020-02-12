package resolver

import (
    "fmt"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
    "piot-server/test"
)

func TestOrgsGet(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateOrg(t, ctx, "org1")

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
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    orgId := test.CreateOrg(t, ctx, "org1")
    userId := test.CreateUser(t, ctx, "org1user@test.com", "")
    test.AddOrgUser(t, ctx, orgId, userId)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            {
                org(id: "%s") {
                    name,
                    users {email},
                    influxdb, influxdb_username, influxdb_password,
                    mysqldb, mysqldb_username, mysqldb_password,
                }
            }
        `, orgId.Hex()),
        ExpectedResult: `
            {
                "org": {
                    "name": "org1",
                    "users": [{"email": "org1user@test.com"}],
                    "influxdb": "db",
                    "influxdb_username": "db-username",
                    "influxdb_password": "db-password",
                    "mysqldb": "mysqldb",
                    "mysqldb_username": "mysqldb-username",
                    "mysqldb_password": "mysqldb-password"
                }
            }
        `,
    })
}

func TestAddOrgUser(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    userId := test.CreateUser(t, ctx, "user1@test.com", "")
    orgId := test.CreateOrg(t, ctx, "test-org")
    org2Id := test.CreateOrg(t, ctx, "test-org2")
    test.CreateOrg(t, ctx, "test-org3")

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
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    userId := test.CreateUser(t, ctx, "user1@test.com", "")
    orgId := test.CreateOrg(t, ctx, "test-org")
    org2Id := test.CreateOrg(t, ctx, "test-org2")
    test.AddOrgUser(t, ctx, orgId, userId)
    test.AddOrgUser(t, ctx, org2Id, userId)

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
