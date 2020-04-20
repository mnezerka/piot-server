package main

import (
    "context"
    "errors"
    "time"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "piot-server/model"
)

type Users struct {
    log *logging.Logger
    db *mongo.Database
}

func NewUsers(log *logging.Logger, db *mongo.Database) *Users{
    return &Users{log: log, db: db}
}

func (t *Users) FindByEmail(email string) (*model.User, error) {
    t.log.Debugf("Get user by email: %s", email)

    var user model.User

    err := t.db.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
    if err != nil {
        t.log.Errorf("Users service error: %v", err)
        return nil, err
    }

    t.log.Debugf("User identified by email <%s> was found.", email)

    // fetch user orgs
    orgs, err := t.FindUserOrgs(user.Id)
    if err != nil {
        t.log.Errorf("Users service error: fetching user orgs failed (%v)", err)
        return nil, err
    }
    user.Orgs = orgs
    return &user, nil
}

func (t *Users) FindUserOrgs(id primitive.ObjectID) ([]model.Org, error) {
    var result []model.Org

    t.log.Debugf("Querying orgs for user: %s", id.Hex())

    collection := t.db.Collection("orgusers")

    // filter orusers to current (single) user
    stage_match := bson.M{"$match": bson.M{"user_id": id}}

    // find orgs details
    stage_lookup := bson.M{"$lookup": bson.M{"from": "orgs", "localField": "org_id", "foreignField": "_id", "as": "orgs"}}

    // expand orgs
    stage_unwind := bson.M{"$unwind": "$orgs"}

    // replace root
    stage_new_root := bson.M{"$replaceWith": "$orgs"}

    pipeline := []bson.M{stage_match, stage_lookup, stage_unwind, stage_new_root}

    cur, err := collection.Aggregate(context.TODO(), pipeline)
    if err != nil {
        t.log.Errorf("Error while querying user orgs: %v", err)
        return result, err
    }
    defer cur.Close(context.TODO())

    for cur.Next(context.TODO()) {
        var org model.Org
        if err := cur.Decode(&org); err != nil {
            t.log.Errorf("Error while querying user orgs: %v", err)
            return result, err
        }
        result = append(result, org)
    }

    if err := cur.Err(); err != nil {
        t.log.Errorf("Error while querying user orgs: %v", err)
        return result, err
    }

    return result, nil
}

func (t *Users) SetActiveOrg(id primitive.ObjectID, orgId primitive.ObjectID) (error) {
    t.log.Debugf("Setting user <%s> active org to to <%s>", id.Hex(), orgId.Hex())

    _, err := t.db.Collection("users").UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"active_org_id": orgId}})
    if err != nil {
        t.log.Errorf("User %s cannot be updated (%v)", id.Hex(), err)
        return errors.New("Error while updating user active org")
    }

    return nil
}

func (t *Users) Create(email, password string) (*model.User, error) {

    // check required attributes
    if len(email) == 0 {
        return nil, errors.New("Email field is empty or not specified!")
    }

    if len(password) == 0 {
        return nil, errors.New("Password field is empty or not specified!")
    }
    if !ValidateEmail(email) {
        return nil, errors.New("Email field has wrong format!")
    }

    // try to find existing user
    var user model.User
    collection := t.db.Collection("users")
    err := collection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
    if err == nil {
        return nil, errors.New("User identified by this email already exists!")
    }

    // generate hash for given password (we don't store passwords in plain form)
    hash, err := GetPasswordHash(password)
    if err != nil {
        return nil, errors.New("Error while hashing password, try again")
    }

    // user does not exist -> create new one
    user.Email = email
    user.Password = hash
    user.Created = int32(time.Now().Unix())

    res, err := collection.InsertOne(context.TODO(), user)
    if err != nil {
        return nil, errors.New("User while creating user, try again")
    }

    user.Id = res.InsertedID.(primitive.ObjectID)

    return &user, nil
}
