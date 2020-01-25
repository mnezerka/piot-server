package resolver

import (
    "errors"
    "strings"
    "time"
    "piot-server/model"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    graphql "github.com/graph-gophers/graphql-go"
    "piot-server/service"
)

type thingUpdateInput struct {
    Id      graphql.ID
    Name    *string
    Alias   *string
    Enabled *bool
    OrgId   *graphql.ID
}

type thingSensorDataUpdateInput struct {
    Id      graphql.ID
    Class *string
    StoreInfluxDb *bool
    MeasurementTopic *string
    MeasurementValue *string
}

type ThingResolver struct {
    ctx context.Context
    t *model.Thing
}

func (r *ThingResolver) Id() graphql.ID {
    return graphql.ID(r.t.Id.Hex())
}

func (r *ThingResolver) Name() string {
    return r.t.Name
}

func (r *ThingResolver) Alias() string {
    return r.t.Alias
}

func (r *ThingResolver) Type() string {
    return r.t.Type
}

func (r *ThingResolver) Enabled() bool {
    return r.t.Enabled
}

func (r *ThingResolver) Created() int32 {
    return r.t.Created
}

func (r *ThingResolver) LastSeen() int32 {
    return r.t.LastSeen
}

func (r *ThingResolver) Org() *OrgResolver {

    r.ctx.Value("log").(*logging.Logger).Debugf("GQL: Fetching org for thing: %s", r.t.Id.Hex())

    if r.t.OrgId != primitive.NilObjectID {

        orgs := r.ctx.Value("orgs").(*service.Orgs)

        org, err := orgs.Get(r.ctx, r.t.OrgId)
        if err != nil {
            r.ctx.Value("log").(*logging.Logger).Errorf("GQL: Fetching org %v for thing %v failed", r.t.OrgId, r.t.Id)
        } else {
            return &OrgResolver{r.ctx, org}
        }
    }

    return nil
}

func (r *ThingResolver) Parent() *ThingResolver {

    r.ctx.Value("log").(*logging.Logger).Debugf("GQL: Fetching parent for thing: %s", r.t.Id.Hex())

    if r.t.ParentId != primitive.NilObjectID {

        things := r.ctx.Value("things").(*service.Things)

        parentThing , err := things.Get(r.ctx, r.t.ParentId)
        if err != nil {
            r.ctx.Value("log").(*logging.Logger).Errorf("GQL: Fetching parent %v for thing %v failed", r.t.ParentId, r.t.Id)
        } else {
            return &ThingResolver{r.ctx, parentThing}
        }
    }

    return nil
}

func (r *ThingResolver) AvailabilityTopic() string {
    return r.t.AvailabilityTopic
}

func (r *ThingResolver) AvailabilityYes() string {
    return r.t.AvailabilityYes
}

func (r *ThingResolver) AvailabilityNo() string {
    return r.t.AvailabilityNo
}

func (r *ThingResolver) Sensor() *SensorResolver {

    if r.t.Type == model.THING_TYPE_SENSOR {
        return &SensorResolver{r.ctx, r.t}
    }

    return nil
}

func (r *ThingResolver) Switch() *SwitchResolver {

    if r.t.Type == model.THING_TYPE_SWITCH {
        return &SwitchResolver{r.ctx, r.t}
    }

    return nil
}


/////////////// Sensor Data Resolver

type SensorResolver struct {
    ctx context.Context
    t *model.Thing
}

func (r *SensorResolver) MeasurementTopic() string {

    return r.t.Sensor.MeasurementTopic

    /*

    // if thing is assigned to org
    if r.t.OrgId != primitive.NilObjectID {

        orgs := r.ctx.Value("orgs").(*service.Orgs)

        org, err := orgs.Get(r.ctx, r.t.OrgId)
        if err != nil {
            r.ctx.Value("log").(*logging.Logger).Errorf("GQL: Fetching org %v for thing %v failed", r.t.OrgId, r.t.Id)
            return ""
        }

        return strings.Join([]string{org.Name, r.t.Name, r.t.Sensor.MeasurementTopic}, "/")

    }
    return ""
    */
}

func (r *SensorResolver) MeasurementValue() string {
    return r.t.Sensor.MeasurementValue
}


func (r *SensorResolver) Value() string {
    return r.t.Sensor.Value
}

func (r *SensorResolver) Unit() string {
    return r.t.Sensor.Unit
}

func (r *SensorResolver) Class() string {
    return r.t.Sensor.Class
}

func (r *SensorResolver) StoreInfluxDb() bool {
    return r.t.Sensor.StoreInfluxDb
}

/////////////// Switch Data Resolver

type SwitchResolver struct {
    ctx context.Context
    t *model.Thing
}

func (r *SwitchResolver) StateTopic() string {

    // if thing is assigned to org
    if r.t.OrgId != primitive.NilObjectID {

        orgs := r.ctx.Value("orgs").(*service.Orgs)

        org, err := orgs.Get(r.ctx, r.t.OrgId)
        if err != nil {
            r.ctx.Value("log").(*logging.Logger).Errorf("GQL: Fetching org %v for thing %v failed", r.t.OrgId, r.t.Id)
            return ""
        }

        return strings.Join([]string{org.Name, r.t.Name, r.t.Switch.StateTopic}, "/")

    }
    return ""
}

func (r *SwitchResolver) StateOn() string {
    return r.t.Switch.StateOn
}

func (r *SwitchResolver) StateOff() string {
    return r.t.Switch.StateOff
}

func (r *SwitchResolver) CommandTopic() string {

    // if thing is assigned to org
    if r.t.OrgId != primitive.NilObjectID {

        orgs := r.ctx.Value("orgs").(*service.Orgs)

        org, err := orgs.Get(r.ctx, r.t.OrgId)
        if err != nil {
            r.ctx.Value("log").(*logging.Logger).Errorf("GQL: Fetching org %v for thing %v failed", r.t.OrgId, r.t.Id)
            return ""
        }

        return strings.Join([]string{org.Name, r.t.Name, r.t.Switch.CommandTopic}, "/")

    }
    return ""
}

func (r *SwitchResolver) CommandOn() string {
    return r.t.Switch.CommandOn
}

func (r *SwitchResolver) CommandOff() string {
    return r.t.Switch.CommandOff
}

func (r *SwitchResolver) StoreInfluxDb() bool {
    return r.t.Switch.StoreInfluxDb
}


/////////////// Resolver

func (r *Resolver) Thing(ctx context.Context, args struct {Id graphql.ID}) (*ThingResolver, error) {

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Fetch thing: %v", args.Id)

    db := ctx.Value("db").(*mongo.Database)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(string(args.Id))
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, errors.New("Cannot decode ID")
    }

    thing := model.Thing{}

    collection := db.Collection("things")
    err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&thing)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Retrieved thing %v", thing)
    return &ThingResolver{ctx, &thing}, nil
}

func (r *Resolver) Things(ctx context.Context) ([]*ThingResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    collection := db.Collection("things")

    count, _ := collection.EstimatedDocumentCount(context.TODO())
    ctx.Value("log").(*logging.Logger).Debugf("GQL: Estimated things count %d", count)

    cur, err := collection.Find(context.TODO(), bson.D{})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
        return nil, err
    }
    defer cur.Close(context.TODO())

    var result []*ThingResolver

    for cur.Next(context.TODO()) {
        // To decode into a struct, use cursor.Decode()
        thing := model.Thing{}
        err := cur.Decode(&thing)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
            return nil, err
        }
        result = append(result, &ThingResolver{ctx, &thing})
    }

    if err := cur.Err(); err != nil {
      return nil, err
    }

    return result, nil
}

func (r *Resolver) CreateThing(ctx context.Context, args *struct {Name string; Type string}) (*ThingResolver, error) {

    thing := &model.Thing{
        Name: args.Name,
        Type: args.Type,
        Created: int32(time.Now().Unix()),
    }

    ctx.Value("log").(*logging.Logger).Infof("Creating thing %s of type %s", args.Name, args.Type)

    if args.Type != model.THING_TYPE_DEVICE && args.Type != model.THING_TYPE_SENSOR && args.Type != model.THING_TYPE_SWITCH {
        return nil, errors.New("Unknown type of Thing")
    }

    db := ctx.Value("db").(*mongo.Database)

    // try to find existing thing
    var existingThing model.Thing
    collection := db.Collection("things")
    err := collection.FindOne(ctx, bson.D{{"name", args.Name}}).Decode(&existingThing)
    if err == nil {
        return nil, errors.New("Thing of such name already exists!")
    }

    // thing does not exist -> create new one
    _, err = collection.InsertOne(ctx, thing)
    if err != nil {
        return nil, errors.New("Error while creating thing")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Created thing: %v", *thing)

    return &ThingResolver{ctx, thing}, nil
}

func (r *Resolver) UpdateThing(ctx context.Context, args struct {Thing thingUpdateInput}) (*ThingResolver, error) {

    ctx.Value("log").(*logging.Logger).Debugf("Updating thing %s", args.Thing.Id)

    db := ctx.Value("db").(*mongo.Database)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(string(args.Thing.Id))
    if err != nil {
        return nil, err
    }

    // try to find thing to be updated
    var thing model.Thing
    collection := db.Collection("things")
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&thing)
    if err != nil {
        return nil, errors.New("Thing does not exist")
    }

    // try to find similar thing matching new name
    if args.Thing.Name != nil {
        var similarThing model.Thing
        err := collection.FindOne(ctx, bson.M{"$and": []bson.M{bson.M{"name": args.Thing.Name}, bson.M{"_id": bson.M{"$ne": id}}}}).Decode(&similarThing)
        if err == nil {
            return nil, errors.New("Thing of such name already exists")
        }
    }

    // thing exists -> update it
    updateFields := bson.M{}
    if args.Thing.Name != nil { updateFields["name"] = *args.Thing.Name}
    if args.Thing.Alias != nil { updateFields["alias"] = *args.Thing.Alias}
    if args.Thing.Enabled != nil { updateFields["enabled"] = *args.Thing.Enabled}
    if args.Thing.OrgId != nil {
        // create ObjectID from string
        orgId, err := primitive.ObjectIDFromHex(string(*args.Thing.OrgId))
        if err != nil {
            return nil, err
        }
        updateFields["org_id"] = orgId
    }
    update := bson.M{"$set": updateFields}

    _, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Updating thing failed %v", err)
        return nil, errors.New("Error while updating thing")
    }

    // read thing
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&thing)
    if err != nil {
        return nil, errors.New("Cannot fetch thing data")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Thing updated %v", thing)
    return &ThingResolver{ctx, &thing}, nil
}

func (r *Resolver) UpdateThingSensorData(ctx context.Context, args struct {Data thingSensorDataUpdateInput}) (*ThingResolver, error) {

    ctx.Value("log").(*logging.Logger).Debugf("Updating thing %s sensor data", args.Data.Id)

    db := ctx.Value("db").(*mongo.Database)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(string(args.Data.Id))
    if err != nil {
        return nil, err
    }

    // try to find thing to be updated
    var thing model.Thing
    collection := db.Collection("things")
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&thing)
    if err != nil {
        return nil, errors.New("Thing does not exist")
    }

    // thing exists -> update it
    updateFields := bson.M{}
    if args.Data.Class != nil { updateFields["sensor.class"] = *args.Data.Class}
    if args.Data.StoreInfluxDb != nil { updateFields["sensor.store_influxdb"] = *args.Data.StoreInfluxDb}
    if args.Data.MeasurementTopic != nil { updateFields["sensor.measurement_topic"] = *args.Data.MeasurementTopic}
    if args.Data.MeasurementValue != nil { updateFields["sensor.measurement_value"] = *args.Data.MeasurementValue}
    update := bson.M{"$set": updateFields}

    _, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Updating thing failed %v", err)
        return nil, errors.New("Error while updating thing")
    }

    // read thing
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&thing)
    if err != nil {
        return nil, errors.New("Cannot fetch thing data")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Thing sensor data updated %v", thing)
    return &ThingResolver{ctx, &thing}, nil
}
