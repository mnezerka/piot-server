package resolver

import (
    "piot-server/model"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"
)

type ThingResolver struct {
    d *model.Thing
}

func (r *ThingResolver) Name() string {
    return r.d.Name
}

func (r *ThingResolver) Type() string {
    return r.d.Type
}

func (r *ThingResolver) Available () bool {
    return r.d.Available
}

func (r *ThingResolver) Created() int32 {
    return r.d.Created
}

func (r *ThingResolver) Org () *OrgResolver {
    // TODO fetch customer
    return nil
}


func (r *Resolver) Thing(ctx context.Context, args struct {Id string}) (*ThingResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    thing := model.Thing{}

    collection := db.Collection("things")
    err := collection.FindOne(context.TODO(), bson.D{{"id", args.Id}}).Decode(&thing)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    return &ThingResolver{&thing}, nil
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
        result = append(result, &ThingResolver{&thing})
    }

    if err := cur.Err(); err != nil {
      return nil, err
    }

    return result, nil
}
