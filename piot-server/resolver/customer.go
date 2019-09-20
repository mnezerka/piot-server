package resolver

import (
    "errors"
    "time"
    "piot-server/model"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

/////////// Customer Resolver

type CustomerResolver struct {
    c *model.Customer
}

func (r *CustomerResolver) Id() string {
    return r.c.Id.Hex()
}

func (r *CustomerResolver) Name() string {
    return r.c.Name
}

func (r *CustomerResolver) Description() string {
    return r.c.Description
}

func (r *CustomerResolver) Created() int32 {
    return r.c.Created
}


/////////// Resolver

func (r *Resolver) Customer(ctx context.Context, args struct {Id string}) (*CustomerResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    customer := model.Customer{}

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(args.Id)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    collection := db.Collection("customers")
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    return &CustomerResolver{&customer}, nil
}

func (r *Resolver) Customers(ctx context.Context) ([]*CustomerResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    collection := db.Collection("customers")

    count, _ := collection.EstimatedDocumentCount(context.TODO())
    ctx.Value("log").(*logging.Logger).Debugf("GQL: Estimated customers count %d", count)

    cur, err := collection.Find(ctx, bson.M{})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
        return nil, err
    }
    defer cur.Close(ctx)

    var result []*CustomerResolver

    for cur.Next(ctx) {
        // To decode into a struct, use cursor.Decode()
        var customer model.Customer
        if err := cur.Decode(&customer); err != nil {
            ctx.Value("log").(*logging.Logger).Debugf("GQL: After decode %v", err)
            ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
            return nil, err
        }
        result = append(result, &CustomerResolver{&customer})
    }

    ctx.Value("log").(*logging.Logger).Debug("Have customers")

    if err := cur.Err(); err != nil {
      return nil, err
    }

    return result, nil
}

func (r *Resolver) CreateCustomer(ctx context.Context, args *struct {Name string; Description string}) (*CustomerResolver, error) {

    customer := &model.Customer{
        Name: args.Name,
        Description: args.Description,
        Created: int32(time.Now().Unix()),
    }

    ctx.Value("log").(*logging.Logger).Infof("Creating customer %s", args.Name)

    db := ctx.Value("db").(*mongo.Database)

    // try to find existing user
    var customerExisting model.Customer
    collection := db.Collection("customers")
    err := collection.FindOne(context.TODO(), bson.D{{"name", args.Name}}).Decode(&customerExisting)
    if err == nil {
        return nil, errors.New("User of such name already exists!")
    }

    // user does not exist -> create new one
    _, err = collection.InsertOne(context.TODO(), customer)
    if err != nil {
        return nil, errors.New("Error while creating customer")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Created customer: %v", *customer)

    return &CustomerResolver{customer}, nil
}

func (r *Resolver) UpdateCustomer(ctx context.Context, args *struct {Id string; Name *string; Description *string}) (*CustomerResolver, error) {

    ctx.Value("log").(*logging.Logger).Debugf("Updating customer %ss", args.Id)

    db := ctx.Value("db").(*mongo.Database)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(args.Id)
    if err != nil {
        return nil, err
    }

    // try to find customer to be updated
    var customer model.Customer
    collection := db.Collection("customers")
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
    if err != nil {
        return nil, errors.New("Customer does not exist")
    }

    // try to find similar customer matching new name
    if args.Name != nil {
        var similarCustomer model.Customer
        collection := db.Collection("customers")
        err := collection.FindOne(ctx, bson.M{"$and": []bson.M{bson.M{"name": args.Name}, bson.M{"_id": bson.M{"$ne": id}}}}).Decode(&similarCustomer)
        if err == nil {
            return nil, errors.New("Customer of such name already exists")
        }
    }

    // customer exists -> update it
    updateFields := bson.M{}
    if args.Name != nil { updateFields["name"] = args.Name}
    if args.Description != nil { updateFields["description"] = args.Description}
    update := bson.M{"$set": updateFields}

    _, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Updating customer failed %v", err)
        return nil, errors.New("Error while updating customer")
    }

    // read customer
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
    if err != nil {
        return nil, errors.New("Cannot fetch customer data")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Customer updated %v", customer)
    return &CustomerResolver{&customer}, nil
}

