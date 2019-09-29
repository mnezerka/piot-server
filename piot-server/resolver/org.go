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
    graphql "github.com/graph-gophers/graphql-go"
)

/////////// Org Resolver

type OrgResolver struct {
    ctx context.Context
    org *model.Org
}

func (r *OrgResolver) Id() graphql.ID {
    return graphql.ID(r.org.Id.Hex())
}

func (r *OrgResolver) Name() string {
    return r.org.Name
}

func (r *OrgResolver) Description() string {
    return r.org.Description
}

func (r *OrgResolver) Created() int32 {
    return r.org.Created
}

    // select all users that are assigned to current org
func (r *OrgResolver) Users() []*UserResolver {

    var result []*UserResolver

    r.ctx.Value("log").(*logging.Logger).Debugf("GQL: Fetching users for org: %s", r.org.Id.Hex())

    db := r.ctx.Value("db").(*mongo.Database)

    collection := db.Collection("orgusers")

    // filter orusers to current (single) org
    stage_match := bson.M{"$match": bson.M{"org_id": r.org.Id}}

    // find assignments to orgs
    stage_lookup := bson.M{"$lookup": bson.M{"from": "users", "localField": "user_id", "foreignField": "_id", "as": "users"}}

    // unwind users
    stage_unwind := bson.M{"$unwind": "$users"}

    // replace root
    stage_new_root := bson.M{"$replaceWith": "$users"}

    pipeline := []bson.M{stage_match, stage_lookup, stage_unwind, stage_new_root}

    //r.ctx.Value("log").(*logging.Logger).Debugf("GQL: Pipeline %v", pipeline)

    cur, err := collection.Aggregate(r.ctx, pipeline)
    if err != nil {
        r.ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
        return result
    }
    defer cur.Close(r.ctx)

    for cur.Next(r.ctx) {
        //r.ctx.Value("log").(*logging.Logger).Debugf("Org users iteration %v", cur.Current)

        var user model.User
        if err := cur.Decode(&user); err != nil {
            r.ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
            return result
        }
        result = append(result, &UserResolver{ctx, &user})
    }

    if err := cur.Err(); err != nil {
        r.ctx.Value("log").(*logging.Logger).Errorf("GQL: error during cursor processing: %v", err)
        return result
    }

    return result
}

/////////// Resolver

func (r *Resolver) Org(ctx context.Context, args struct {Id graphql.ID}) (*OrgResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    org := model.Org{}

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Fetching org %v", args.Id)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(string(args.Id))
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, errors.New("Cannot decode ID")
    }

    collection := db.Collection("orgs")
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&org)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    return &OrgResolver{ctx, &org}, nil
}

func (r *Resolver) Orgs(ctx context.Context) ([]*OrgResolver, error) {

    db := ctx.Value("db").(*mongo.Database)

    collection := db.Collection("orgs")

    count, _ := collection.EstimatedDocumentCount(context.TODO())
    ctx.Value("log").(*logging.Logger).Debugf("GQL: Estimated orgs count %d", count)

    cur, err := collection.Find(ctx, bson.M{})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
        return nil, err
    }
    defer cur.Close(ctx)

    var result []*OrgResolver

    for cur.Next(ctx) {
        // To decode into a struct, use cursor.Decode()
        var org model.Org
        if err := cur.Decode(&org); err != nil {
            ctx.Value("log").(*logging.Logger).Debugf("GQL: After decode %v", err)
            ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
            return nil, err
        }
        result = append(result, &OrgResolver{ctx, &org})
    }

    ctx.Value("log").(*logging.Logger).Debug("Have orgs")

    if err := cur.Err(); err != nil {
      return nil, err
    }

    return result, nil
}

func (r *Resolver) CreateOrg(ctx context.Context, args *struct {Name string; Description string}) (*OrgResolver, error) {

    org := &model.Org{
        Name: args.Name,
        Description: args.Description,
        Created: int32(time.Now().Unix()),
    }

    ctx.Value("log").(*logging.Logger).Infof("Creating org %s", args.Name)

    db := ctx.Value("db").(*mongo.Database)

    // try to find existing user
    var orgExisting model.Org
    collection := db.Collection("orgs")
    err := collection.FindOne(ctx, bson.D{{"name", args.Name}}).Decode(&orgExisting)
    if err == nil {
        return nil, errors.New("Organization of such name already exists!")
    }

    // org does not exist -> create new one
    _, err = collection.InsertOne(ctx, org)
    if err != nil {
        return nil, errors.New("Error while creating organizaton")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Created organization: %v", *org)

    return &OrgResolver{ctx, org}, nil
}

func (r *Resolver) UpdateOrg(ctx context.Context, args *struct {Id string; Name *string; Description *string}) (*OrgResolver, error) {

    ctx.Value("log").(*logging.Logger).Debugf("Updating org %ss", args.Id)

    db := ctx.Value("db").(*mongo.Database)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(args.Id)
    if err != nil {
        return nil, err
    }

    // try to find org to be updated
    var org model.Org
    collection := db.Collection("orgs")
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&org)
    if err != nil {
        return nil, errors.New("Org does not exist")
    }

    // try to find similar org matching new name
    if args.Name != nil {
        var similarOrg model.Org
        collection := db.Collection("orgs")
        err := collection.FindOne(ctx, bson.M{"$and": []bson.M{bson.M{"name": args.Name}, bson.M{"_id": bson.M{"$ne": id}}}}).Decode(&similarOrg)
        if err == nil {
            return nil, errors.New("Org of such name already exists")
        }
    }

    // org exists -> update it
    updateFields := bson.M{}
    if args.Name != nil { updateFields["name"] = args.Name}
    if args.Description != nil { updateFields["description"] = args.Description}
    update := bson.M{"$set": updateFields}

    _, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Updating org failed %v", err)
        return nil, errors.New("Error while updating org")
    }

    // read org
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&org)
    if err != nil {
        return nil, errors.New("Cannot fetch org data")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Org updated %v", org)
    return &OrgResolver{ctx, &org}, nil
}

func (r *Resolver) AssignOrgUser(ctx context.Context, args *struct {OrgId graphql.ID; UserId graphql.ID}) (*bool, error) {

    ctx.Value("log").(*logging.Logger).Debugf("Assigning user %s to org %s", args.UserId, args.OrgId)

    db := ctx.Value("db").(*mongo.Database)

    // create ObjectIDs from string
    orgId, err := primitive.ObjectIDFromHex(string(args.OrgId))
    if err != nil {
        return nil, err
    }
    userId, err := primitive.ObjectIDFromHex(string(args.UserId))
    if err != nil {
        return nil, err
    }

    // try to find existing assignment
    var similarOrgUser model.OrgUser
    collection := db.Collection("orgusers")
    err = collection.FindOne(ctx, bson.M{"$and": []bson.M{bson.M{"user_id": userId}, bson.M{"org_id": bson.M{"$ne": orgId}}}}).Decode(&similarOrgUser)
    if err == nil {
        return nil, errors.New("User is allready assigned to given organization")
    }

    // assignment does not exist -> create new one
    orgUser := &model.OrgUser{
        UserId: userId,
        OrgId: orgId,
        Created: int32(time.Now().Unix()),
    }
    _, err = collection.InsertOne(ctx, orgUser)
    if err != nil {
        return nil, errors.New("Error while assigning user to organization")
    }

    ctx.Value("log").(*logging.Logger).Debugf("User %s assigned to Org %s", userId, orgId)
    return nil, nil
}
