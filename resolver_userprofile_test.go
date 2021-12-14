package main_test

import (
	"fmt"
	"piot-server/schema"
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/gqltesting"
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
		Context: AuthContext(t, userId, orgId),
		Schema:  schema,
		Query: `
            {
                userProfile() {email, org_id, is_admin, orgs {name}}
            }
        `,
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

func TestUserProfileUpdate(t *testing.T) {
	db := GetDb(t)
	CleanDb(t, db)
	orgId := CreateOrg(t, db, "org1")
	org2Id := CreateOrg(t, db, "org2")
	userId := CreateUser(t, db, "user1@test.com", "")
	AddOrgUser(t, db, orgId, userId)
	AddOrgUser(t, db, org2Id, userId)

	schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

	// set second org as active
	gqltesting.RunTest(t, &gqltesting.Test{
		Context: AuthContext(t, userId, orgId),
		Schema:  schema,
		Query: fmt.Sprintf(`
            mutation {
                updateUserProfile(profile: {org_id: "%s"}) {org_id}
            }
        `, org2Id.Hex()),
		ExpectedResult: fmt.Sprintf(`
            {
                "updateUserProfile": {
                    "org_id": "%s"
                }
            }
        `, org2Id.Hex()),
	})

	// set first org as active
	gqltesting.RunTest(t, &gqltesting.Test{
		Context: AuthContext(t, userId, orgId),
		Schema:  schema,
		Query: fmt.Sprintf(`
            mutation {
                updateUserProfile(profile: {org_id: "%s"}) {org_id}
            }
        `, orgId.Hex()),
		ExpectedResult: fmt.Sprintf(`
            {
                "updateUserProfile": {
                    "org_id": "%s"
                }
            }
        `, orgId.Hex()),
	})
}
