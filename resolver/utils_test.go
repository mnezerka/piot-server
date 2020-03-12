package resolver_test

import(
    "testing"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/mnezerka/go-piot/test"
    "piot-server/resolver"
)

func getResolver(t *testing.T, db *mongo.Database) *resolver.Resolver{
    log := test.GetLogger(t)
    users := test.GetUsers(t, log, db)
    orgs := test.GetOrgs(t, log, db)
    things := test.GetThings(t, log, db)

    return resolver.NewResolver(log, db, orgs, users, things)
}
