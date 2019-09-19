package resolver

import (
    "errors"
    "time"
    "piot-server/model"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"
)

type CustomerResolver struct {
    c *model.Customer
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

func (r *Resolver) Customer(ctx context.Context, args struct {Name string}) (*CustomerResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    customer := model.Customer{}

    collection := db.Collection("customers")
    err := collection.FindOne(context.TODO(), bson.D{{"name", args.Name}}).Decode(&customer)
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

    cur, err := collection.Find(context.TODO(), bson.D{})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
        return nil, err
    }
    defer cur.Close(context.TODO())

    var result []*CustomerResolver

    for cur.Next(context.TODO()) {
        // To decode into a struct, use cursor.Decode()
        customer := model.Customer{}
        err := cur.Decode(&customer)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
            return nil, err
        }
        result = append(result, &CustomerResolver{&customer})
    }

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

func (r *Resolver) UpdateCustomer(ctx context.Context, args *struct {Name string; NewName *string; NewDescription *string}) (*CustomerResolver, error) {

    ctx.Value("log").(*logging.Logger).Infof("Updating customer %s", args.Name)

    db := ctx.Value("db").(*mongo.Database)

    // try to find customer to be updated
    var customer model.Customer
    collection := db.Collection("customers")
    err := collection.FindOne(context.TODO(), bson.D{{"name", args.Name}}).Decode(&customer)
    if err != nil {
        return nil, errors.New("Customer does not exist")
    }

    // try to find similar customer matching new name
    if args.NewName != nil {
        var similarCustomer model.Customer
        collection := db.Collection("customers")
        err := collection.FindOne(context.TODO(), bson.D{{"name", args.NewName}}).Decode(&similarCustomer)
        ctx.Value("log").(*logging.Logger).Infof("Similar customer search result %v %v", err, similarCustomer)
    }
    if err == nil {
        return nil, errors.New("Customer of such name already exists")
    }


    // customer exists -> update it
    /*
    _, err = collection.InsertOne(context.TODO(), customer)
    if err != nil {
        return nil, errors.New("Error while creating customer")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Created customer: %v", *customer)
    */
    return &CustomerResolver{&customer}, nil
}

