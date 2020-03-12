package resolver

import (
    "errors"
    "time"
    "github.com/mnezerka/go-piot"
    "github.com/mnezerka/go-piot/model"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    graphql "github.com/graph-gophers/graphql-go"
)

type userUpdateInput struct {
    Id      graphql.ID
    Email   *string
    OrgId   *graphql.ID
}

type userCreateInput struct {
    Email   string
    OrgId   *graphql.ID
}

/////////// User Resolver

type UserResolver struct {
    log *logging.Logger
    users *piot.Users
    db *mongo.Database
    u *model.User
}

func (r *UserResolver) Id() graphql.ID {
    return graphql.ID(r.u.Id.Hex())
}

func (r *UserResolver) Email() string {
    return r.u.Email
}

func (r *UserResolver) Password() string {
    maskedPassword := "********"
    return maskedPassword
}

func (r *UserResolver) Created() int32 {
    return r.u.Created
}

func (r *UserResolver) Orgs() []*OrgResolver {
    var result []*OrgResolver

    // get all orgs assigned to user
    orgs, err := r.users.FindUserOrgs(r.u.Id)
    if err != nil {
        r.log.Errorf("GQL: error : %v", err)
        return result
    }

    // convert orgs to org resolvers
    for i := 0; i < len(orgs); i++ {
        result = append(result, &OrgResolver{r.log, r.db, r.users, &orgs[i]})
    }

    return result
}

/////////// UserProfileResolver

type UserProfileResolver struct {
    u *model.User
}

func (r *UserProfileResolver) Email() string {
    return r.u.Email
}

/////////// Resolver

// get user by email query
func (r *Resolver) User(args struct {Id graphql.ID}) (*UserResolver, error) {

    r.log.Debugf("GQL: Fetch user: %v", args.Id)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(string(args.Id))
    if err != nil {
        r.log.Errorf("Graphql error : %v", err)
        return nil, errors.New("Cannot decode ID")
    }

    user := model.User{}

    collection := r.db.Collection("users")
    err = collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
    if err != nil {
        r.log.Errorf("Graphql error : %v", err)
        return nil, err
    }

    r.log.Debugf("GQL: Retrieved user %v", user)
    return &UserResolver{r.log, r.users, r.db, &user}, nil
}

// get users query
func (r *Resolver) Users(ctx context.Context) ([]*UserResolver, error) {

    collection := r.db.Collection("users")

    count, _ := collection.EstimatedDocumentCount(context.TODO())
    r.log.Debugf("GQL: Estimated users count %d", count)

    cur, err := collection.Find(context.TODO(), bson.D{})
    if err != nil {
        r.log.Errorf("GQL: error : %v", err)
        return nil, err
    }
    defer cur.Close(context.TODO())

    var result []*UserResolver

    for cur.Next(context.TODO()) {
        // To decode into a struct, use cursor.Decode()
        user := model.User{}
        err := cur.Decode(&user)
        if err != nil {
            r.log.Errorf("GQL: error : %v", err)
            return nil, err
        }
        result = append(result, &UserResolver{r.log, r.users, r.db, &user})
    }

    if err := cur.Err(); err != nil {
      return nil, err
    }

    return result, nil
}

// get active user profile
func (r *Resolver) UserProfile(ctx context.Context) (*UserProfileResolver, error) {

    currentUserEmail := ctx.Value("user_email").(*string)

    r.log.Debugf("GQL: Getting user profile for %s", *currentUserEmail)

    user := model.User{}

    collection := r.db.Collection("users")
    err := collection.FindOne(context.TODO(), bson.D{{"email", currentUserEmail}}).Decode(&user)
    if err != nil {
        r.log.Errorf("Graphql error : %v", err)
        return nil, err
    }

    return &UserProfileResolver{&user}, nil
}

func (r *Resolver) CreateUser(ctx context.Context, args struct {User userCreateInput}) (*UserResolver, error) {

    user := &model.User{
        Email: args.User.Email,
        Created: int32(time.Now().Unix()),
    }

    r.log.Infof("Creating user %s", args.User.Email)

    // try to find existing user of same email
    var userExisting model.User
    collection := r.db.Collection("users")
    err := collection.FindOne(context.TODO(), bson.D{{"email", args.User.Email}}).Decode(&userExisting)
    if err == nil {
        return nil, errors.New("User of such email already exists!")
    }

    // user does not exist -> create new one
    _, err = collection.InsertOne(context.TODO(), user)
    if err != nil {
        return nil, errors.New("Error while creating user")
    }

    r.log.Debugf("Created user: %s", args.User.Email)

    return &UserResolver{r.log, r.users, r.db, user}, nil
}

func (r *Resolver) UpdateUser(args struct {User userUpdateInput}) (*UserResolver, error) {

    r.log.Debugf("Updating user %s", args.User.Id)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(string(args.User.Id))
    if err != nil {
        return nil, err
    }

    // try to find user to be updated
    var user model.User
    collection := r.db.Collection("users")
    err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, errors.New("User does not exist")
    }

    // try to find similar user matching new email
    if args.User.Email != nil {
        var similarUser model.User
        err := collection.FindOne(context.TODO(), bson.M{"$and": []bson.M{bson.M{"email": args.User.Email}, bson.M{"_id": bson.M{"$ne": id}}}}).Decode(&similarUser)
        if err == nil {
            return nil, errors.New("User of such name already exists")
        }
    }

    // user exists -> update it
    updateFields := bson.M{}
    if args.User.Email != nil { updateFields["email"] = args.User.Email}
    update := bson.M{"$set": updateFields}

    _, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
    if err != nil {
        r.log.Errorf("Updating user failed %v", err)
        return nil, errors.New("Error while updating user")
    }

    // read user
    err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, errors.New("Cannot fetch user data")
    }

    r.log.Debugf("User updated %v", user)
    return &UserResolver{r.log, r.users, r.db, &user}, nil
}
