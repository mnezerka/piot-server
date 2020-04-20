package main_test

import (
    "fmt"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
)

func TestUserProfile(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    orgId := CreateOrg(t, db, "org1")
    org2Id := CreateOrg(t, db, "org2")
    userId := CreateUser(t, db, "user1@test.com", "")
    AddOrgUser(t, db, orgId, userId)
    AddOrgUser(t, db, org2Id, userId)

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: AuthContext(t, userId, orgId, false),
        Schema:  schema,
        Query: fmt.Sprintf(`
            {
                userProfile() {email, org_id, is_admin, orgs {name}}
            }
        `),
        ExpectedResult: fmt.Sprintf(`
            {
                "userProfile": {
                    "email": "user1@test.com",
                    "org_id": "%s",
                    "is_admin": false,
                    "orgs": [{"name": "org1"}, {"name": "org2"}]
                }
            }
        `, orgId.Hex()),
    })
}
