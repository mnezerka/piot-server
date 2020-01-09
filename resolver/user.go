package resolver

import (
    "errors"
    "time"
    "piot-server/model"
    "piot-server/service"
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
    ctx context.Context
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

    // get instance of the Users service
    users := r.ctx.Value("users").(*service.Users)

    // get all orgs assigned to user
    orgs, err := users.FindUserOrgs(r.ctx, r.u.Id)
    if err != nil {
        r.ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
        return result
    }

    // convert orgs to org resolvers
    for _, userOrg := range orgs {
        result = append(result, &OrgResolver{r.ctx, &userOrg})
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
func (r *Resolver) User(ctx context.Context, args struct {Id graphql.ID}) (*UserResolver, error) {

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Fetch user: %v", args.Id)

    db := ctx.Value("db").(*mongo.Database)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(string(args.Id))
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, errors.New("Cannot decode ID")
    }

    user := model.User{}

    collection := db.Collection("users")
    err = collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Retrieved user %v", user)
    return &UserResolver{ctx, &user}, nil
}

// get users query
func (r *Resolver) Users(ctx context.Context) ([]*UserResolver, error) {

    currentUserEmail := ctx.Value("user_email").(*string)

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Retrieved users by %s: ", *currentUserEmail)

    db := ctx.Value("db").(*mongo.Database)

    collection := db.Collection("users")

    count, _ := collection.EstimatedDocumentCount(context.TODO())
    ctx.Value("log").(*logging.Logger).Debugf("GQL: Estimated users count %d", count)

    cur, err := collection.Find(context.TODO(), bson.D{})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
        return nil, err
    }
    defer cur.Close(context.TODO())

    var result []*UserResolver

    for cur.Next(context.TODO()) {
        // To decode into a struct, use cursor.Decode()
        user := model.User{}
        err := cur.Decode(&user)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("GQL: error : %v", err)
            return nil, err
        }
        result = append(result, &UserResolver{ctx, &user})
    }

    if err := cur.Err(); err != nil {
      return nil, err
    }

    return result, nil
}

// get active user profile
func (r *Resolver) UserProfile(ctx context.Context) (*UserProfileResolver, error) {

    currentUserEmail := ctx.Value("user_email").(*string)

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Getting user profile for %s", *currentUserEmail)

    db := ctx.Value("db").(*mongo.Database)

    user := model.User{}

    collection := db.Collection("users")
    err := collection.FindOne(context.TODO(), bson.D{{"email", currentUserEmail}}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    return &UserProfileResolver{&user}, nil
}

func (r *Resolver) CreateUser(ctx context.Context, args struct {User userCreateInput}) (*UserResolver, error) {

    user := &model.User{
        Email: args.User.Email,
        Created: int32(time.Now().Unix()),
    }

    ctx.Value("log").(*logging.Logger).Infof("Creating user %s", args.User.Email)

    db := ctx.Value("db").(*mongo.Database)

    // try to find existing user of same email
    var userExisting model.User
    collection := db.Collection("users")
    err := collection.FindOne(ctx, bson.D{{"email", args.User.Email}}).Decode(&userExisting)
    if err == nil {
        return nil, errors.New("User of such email already exists!")
    }

    // user does not exist -> create new one
    _, err = collection.InsertOne(ctx, user)
    if err != nil {
        return nil, errors.New("Error while creating user")
    }

    ctx.Value("log").(*logging.Logger).Debugf("Created user: %s", args.User.Email)

    return &UserResolver{ctx, user}, nil
}

func (r *Resolver) UpdateUser(ctx context.Context, args struct {User userUpdateInput}) (*UserResolver, error) {

    ctx.Value("log").(*logging.Logger).Debugf("Updating user %s", args.User.Id)

    db := ctx.Value("db").(*mongo.Database)

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex(string(args.User.Id))
    if err != nil {
        return nil, err
    }

    // try to find user to be updated
    var user model.User
    collection := db.Collection("users")
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, errors.New("User does not exist")
    }

    // try to find similar user matching new email
    if args.User.Email != nil {
        var similarUser model.User
        err := collection.FindOne(ctx, bson.M{"$and": []bson.M{bson.M{"email": args.User.Email}, bson.M{"_id": bson.M{"$ne": id}}}}).Decode(&similarUser)
        if err == nil {
            return nil, errors.New("User of such name already exists")
        }
    }

    // user exists -> update it
    updateFields := bson.M{}
    if args.User.Email != nil { updateFields["email"] = args.User.Email}
    update := bson.M{"$set": updateFields}

    _, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Updating user failed %v", err)
        return nil, errors.New("Error while updating user")
    }

    // read user
    err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, errors.New("Cannot fetch user data")
    }

    ctx.Value("log").(*logging.Logger).Debugf("User updated %v", user)
    return &UserResolver{ctx, &user}, nil
}
