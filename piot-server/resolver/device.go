package resolver

import (
    "piot-server/model"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"
)

type DeviceResolver struct {
    d *model.Device
}

func (r *DeviceResolver) Name() string {
    return r.d.Name
}

func (r *DeviceResolver) Type() string {
    return r.d.Type
}

func (r *DeviceResolver) Available () bool {
    return r.d.Available
}

func (r *DeviceResolver) Created() int32 {
    return r.d.Created
}

func (r *DeviceResolver) Customer () *CustomerResolver {
    // TODO fetch customer
    return nil
}


func (r *Resolver) Device(ctx context.Context, args struct {Id string}) (*DeviceResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    device := model.Device{}

    collection := db.Collection("devices")
    err := collection.FindOne(context.TODO(), bson.D{{"id", args.Id}}).Decode(&device)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    return &DeviceResolver{&device}, nil
}

func (r *Resolver) Devices(ctx context.Context) ([]*DeviceResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    collection := db.Collection("devices")

    count, _ := collection.EstimatedDocumentCount(context.TODO())
    ctx.Value("log").(*logging.Logger).Debugf("GQL: Estimated devices count %d", count)

    cur, err := collection.Find(context.TODO(), bson.D{})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
        return nil, err
    }
    defer cur.Close(context.TODO())

    var result []*DeviceResolver

    for cur.Next(context.TODO()) {
        // To decode into a struct, use cursor.Decode()
        device := model.Device{}
        err := cur.Decode(&device)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
            return nil, err
        }
        result = append(result, &DeviceResolver{&device})
    }

    if err := cur.Err(); err != nil {
      return nil, err
    }

    return result, nil
}
