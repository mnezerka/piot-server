package service

import (
    "context"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/mnezerka/go-piot/model"
    //"piot-server/model"
)

type Users struct { }

func (t *Users) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    ctx.Value("log").(*logging.Logger).Debugf("Get user by email: %s", email)

    db := ctx.Value("db").(*mongo.Database)

    var user model.User

    err := db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Users service error: %v", err)
        return nil, err
    }

    ctx.Value("log").(*logging.Logger).Debugf("User identified by email <%s> was found.", email)

    // fetch user orgs

    orgs, err := t.FindUserOrgs(ctx, user.Id)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Users service error: fetching user orgs failed (%v)", err)
        return nil, err
    }
    user.Orgs = orgs
    return &user, nil
}

func (t *Users) FindUserOrgs(ctx context.Context, id primitive.ObjectID) ([]model.Org, error) {

    var result []model.Org

    ctx.Value("log").(*logging.Logger).Debugf("Querying orgs for user: %s", id.Hex())

    db := ctx.Value("db").(*mongo.Database)

    collection := db.Collection("orgusers")

    // filter orusers to current (single) user
    stage_match := bson.M{"$match": bson.M{"user_id": id}}

    // find orgs details
    stage_lookup := bson.M{"$lookup": bson.M{"from": "orgs", "localField": "org_id", "foreignField": "_id", "as": "orgs"}}

    // expand orgs
    stage_unwind := bson.M{"$unwind": "$orgs"}

    // replace root
    stage_new_root := bson.M{"$replaceWith": "$orgs"}

    pipeline := []bson.M{stage_match, stage_lookup, stage_unwind, stage_new_root}

    cur, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Error while querying user orgs: %v", err)
        return result, err
    }
    defer cur.Close(ctx)

    for cur.Next(ctx) {
        var org model.Org
        if err := cur.Decode(&org); err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("Error while querying user orgs: %v", err)
            return result, err
        }
        result = append(result, org)
    }

    if err := cur.Err(); err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Error while querying user orgs: %v", err)
        return result, err
    }

    return result, nil
}
