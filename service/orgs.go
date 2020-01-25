package service

import (
    "context"
    "errors"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "piot-server/model"
)

type Orgs struct { }

func (t *Orgs) Get(ctx context.Context, id primitive.ObjectID) (*model.Org, error) {
    ctx.Value("log").(*logging.Logger).Debugf("Get org: %s", id.Hex())

    db := ctx.Value("db").(*mongo.Database)

    var org model.Org

    collection := db.Collection("orgs")
    err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&org)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Org service error : %v", err)
        return nil, err
    }

    return &org, nil
}

func (t *Orgs) GetByName(ctx context.Context, name string) (*model.Org, error) {
    ctx.Value("log").(*logging.Logger).Debugf("Finding org by name <%s>", name)

    db := ctx.Value("db").(*mongo.Database)

    var org model.Org

    // try to find thing in DB by its name
    err := db.Collection("orgs").FindOne(ctx, bson.M{"name": name}).Decode(&org)
    if err != nil {
        //ctx.Value("log").(*logging.Logger).Errorf("Thing %s not found (%v)", name, err)
        return nil, errors.New("Org not found")
    }

    return &org, nil
}

