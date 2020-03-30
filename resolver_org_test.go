package main_test

import (
    "context"
    "fmt"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
)

func TestOrgsGet(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    CreateOrg(t, db, "org1")
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: context.TODO(),
            Schema: schema,
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
    db := GetDb(t)
    CleanDb(t, db)
    orgId := CreateOrg(t, db, "org1")
    userId := CreateUser(t, db, "org1user@test.com", "")
    AddOrgUser(t, db, orgId, userId)
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
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
    db := GetDb(t)
    CleanDb(t, db)
    userId := CreateUser(t, db, "user1@test.com", "")
    orgId := CreateOrg(t, db, "test-org")
    org2Id := CreateOrg(t, db, "test-org2")
    CreateOrg(t, db, "test-org3")
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    t.Logf("Test adding user %s to org %s", userId, orgId)

    // assign user to the first organization
    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
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
        Context: context.TODO(),
        Schema: schema,
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
    db := GetDb(t)
    CleanDb(t, db)
    userId := CreateUser(t, db, "user1@test.com", "")
    orgId := CreateOrg(t, db, "test-org")
    org2Id := CreateOrg(t, db, "test-org2")
    AddOrgUser(t, db, orgId, userId)
    AddOrgUser(t, db, org2Id, userId)
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    t.Logf("Test remove user %s from org %s", userId, orgId)

    // assign user to the first organization
    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
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
