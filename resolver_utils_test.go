package main_test

import(
    "testing"
    "go.mongodb.org/mongo-driver/mongo"
    "piot-server"
)

func getResolver(t *testing.T, db *mongo.Database) *main.Resolver{
    log := GetLogger(t)
    users := GetUsers(t, log, db)
    orgs := GetOrgs(t, log, db)
    things := GetThings(t, log, db)

    return main.NewResolver(log, db, orgs, users, things)
}
