package resolver

import(
    "github.com/mnezerka/go-piot"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
)

type Resolver struct{
    log *logging.Logger
    db *mongo.Database
    orgs *piot.Orgs
    things *piot.Things
    users *piot.Users
}

func NewResolver(log *logging.Logger, db *mongo.Database, orgs *piot.Orgs, users *piot.Users, things *piot.Things) *Resolver {
    return &Resolver{log: log, db: db, orgs: orgs, things: things, users: users}
}
