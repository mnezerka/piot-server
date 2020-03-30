package main

import (
    "context"
    "errors"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "piot-server/model"
)

type Orgs struct {
    log *logging.Logger
    db *mongo.Database
}

func NewOrgs(log *logging.Logger, db *mongo.Database) *Orgs {
    return &Orgs{log: log, db: db}
}

func (t *Orgs) Get(id primitive.ObjectID) (*model.Org, error) {
    t.log.Debugf("Get org: %s", id.Hex())

    var org model.Org

    collection := t.db.Collection("orgs")
    err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&org)
    if err != nil {
        t.log.Errorf("Org service error : %v", err)
        return nil, err
    }

    return &org, nil
}

func (t *Orgs) GetByName(name string) (*model.Org, error) {
    t.log.Debugf("Finding org by name <%s>", name)

    var org model.Org

    // try to find thing in DB by its name
    err := t.db.Collection("orgs").FindOne(context.TODO(), bson.M{"name": name}).Decode(&org)
    if err != nil {
        return nil, errors.New("Org not found")
    }

    return &org, nil
}
