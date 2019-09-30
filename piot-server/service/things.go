package service

import (
    "context"
    "errors"
    "fmt"
    "time"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "piot-server/model"
)

type ThingsReader interface {
    Find(name *string) (*model.Thing, error)
}

type ThingsWriter interface {
}

type ThingsService interface {
    ThingsReader
    ThingsWriter
}

type Things struct { }

func (t *Things) Find(ctx context.Context, name string) (*model.Thing, error) {
    ctx.Value("log").(*logging.Logger).Debugf("[TH] Find thing: %s", name)

    db := ctx.Value("db").(*mongo.Database)

    var thing model.Thing

    // try to find thing in DB by its name
    err := db.Collection("things").FindOne(ctx, bson.M{"name": name}).Decode(&thing)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("[TH] Thing %s not found (%v)", name, err)
        return nil, errors.New("Thing not found")
    }

    return &thing, nil
}

func (t *Things) Register(ctx context.Context, name string, deviceType string) (*model.Thing, error) {
    ctx.Value("log").(*logging.Logger).Debugf("[TH] Registering new thing: %s of type %s", name, deviceType)

    // check if string of same name already exists
    _, err := t.Find(ctx, name)
    if err == nil {
        return nil, errors.New(fmt.Sprintf("Thing %s already exists", name))
    }

    // thing does not exist -> create new one
    db := ctx.Value("db").(*mongo.Database)

    var thing model.Thing
    thing.Name = name
    thing.Type = deviceType
    thing.Created = int32(time.Now().Unix())

    res, err := db.Collection("things").InsertOne(ctx, thing)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("[TH] Thing %s cannot be stored (%v)", name, err)
        return nil, errors.New("Error while storing new thing")
    }

    thing.Id = res.InsertedID.(primitive.ObjectID)

    return &thing, nil
}

