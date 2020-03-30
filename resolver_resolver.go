package main

import(
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
)

type Resolver struct{
    log *logging.Logger
    db *mongo.Database
    orgs *Orgs
    things *Things
    users *Users
}

func NewResolver(log *logging.Logger, db *mongo.Database, orgs *Orgs, users *Users, things *Things) *Resolver {
    return &Resolver{log: log, db: db, orgs: orgs, things: things, users: users}
}
