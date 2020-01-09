package service

/* Note: checking of res.ModifiedCount could be tricky since it is
 zero when method is called multiple times in one second - last_seen
 attribute is updated only in the first call

if res.ModifiedCount == 0 {
    return fmt.Errorf("Thing <%s> not found", id.Hex())
}
*/

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

type Things struct { }

func (t *Things) Get(ctx context.Context, id primitive.ObjectID) (*model.Thing, error) {
    ctx.Value("log").(*logging.Logger).Debugf("Get thing: %s", id.Hex())

    db := ctx.Value("db").(*mongo.Database)

    var thing model.Thing

    collection := db.Collection("things")
    err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&thing)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Things service error : %v", err)
        return nil, err
    }

    return &thing, nil
}

func (t *Things) Find(ctx context.Context, name string) (*model.Thing, error) {
    ctx.Value("log").(*logging.Logger).Debugf("Finding thing by name <%s>", name)

    db := ctx.Value("db").(*mongo.Database)

    var thing model.Thing

    // try to find thing in DB by its name
    err := db.Collection("things").FindOne(ctx, bson.M{"name": name}).Decode(&thing)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s not found (%v)", name, err)
        return nil, errors.New("Thing not found")
    }

    return &thing, nil
}

func (t *Things) Register(ctx context.Context, name string, deviceType string) (*model.Thing, error) {
    ctx.Value("log").(*logging.Logger).Debugf("Registering new thing: %s of type %s", name, deviceType) 
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
    thing.LastSeen = int32(time.Now().Unix())

    res, err := db.Collection("things").InsertOne(ctx, thing)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s cannot be stored (%v)", name, err)
        return nil, errors.New("Error while storing new thing")
    }

    thing.Id = res.InsertedID.(primitive.ObjectID)

    return &thing, nil
}

func (t *Things) SetParent(ctx context.Context, id primitive.ObjectID, id_parent primitive.ObjectID) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Setting thing <%v>, setting parent to <%v>", id, id_parent)

    db := ctx.Value("db").(*mongo.Database)

    _, err := t.Get(ctx, id)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s not found", id.Hex())
        return errors.New("Child thing not found when setting new parent")
    }

    _, err = t.Get(ctx, id_parent)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s not found", id_parent.Hex())
        return errors.New("Parent thing not found when setting new parent for thing")
    }

    _, err = db.Collection("things").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"parent_id": id_parent}})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s cannot be updated (%v)", id.Hex(), err)
        return errors.New("Error while updating thing parent")
    }

    return nil
}

func (t *Things) SetAvailabilityTopic(ctx context.Context, name string, topic string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Setting thing <%s>, setting avalibility topic to <%s>", name, topic)

    db := ctx.Value("db").(*mongo.Database)

    _, err := db.Collection("things").UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"availability_topic": topic}})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s cannot be updated (%v)", name, err)
        return errors.New("Error while updating thing attributes")
    }

    return nil
}

func (t *Things) SetAvailabilityYesNo(ctx context.Context, name , yes, no string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Setting thing <%s>, setting avalibility topic values to <%s> and <%s>", name, yes, no)

    db := ctx.Value("db").(*mongo.Database)

    _, err := db.Collection("things").UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"availability_yes": yes, "availability_no": no}})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s cannot be updated (%v)", name, err)
        return errors.New("Error while updating thing attributes")
    }

    return nil
}

func (t *Things) SetSensorMeasurementTopic(ctx context.Context, name string, topic string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Setting thing <%s> sensor measurement topic to <%s>", name, topic)

    db := ctx.Value("db").(*mongo.Database)

    _, err := db.Collection("things").UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"sensor.measurement_topic": topic}})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s cannot be updated (%v)", name, err)
        return errors.New("Error while updating thing attributes")
    }

    return nil
}

func (t *Things) SetSensorClass(ctx context.Context, name string, class string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Setting thing <%s> sensor class to <%s>", name, class)

    db := ctx.Value("db").(*mongo.Database)

    _, err := db.Collection("things").UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"sensor.class": class}})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Thing %s cannot be updated (%v)", name, err)
        return errors.New("Error while updating thing attributes")
    }

    return nil
}

func (t *Things) TouchThing(ctx context.Context, id primitive.ObjectID) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Touch thing <%s>", id.Hex())

    db := ctx.Value("db").(*mongo.Database)

    _, err := db.Collection("things").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"last_seen": int32(time.Now().Unix())}})
    if err != nil {
        e := fmt.Errorf("Thing <%s> cannot be touched (%v)", id.Hex(), err)
        ctx.Value("log").(*logging.Logger).Errorf(e.Error())
        return e
    }

    return nil
}
