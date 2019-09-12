package resolver

import (
    "piot-server/model"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"
)

type UserResolver struct {
    u *model.User
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

// get user by email query
func (r *Resolver) User(ctx context.Context, args struct {Email string}) (*UserResolver, error) {

    currentUserEmail := ctx.Value("user_email").(*string)

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Creating User resolver for: %s, triggered by %s", args.Email, *currentUserEmail)

    db := ctx.Value("db").(*mongo.Database)

    user := model.User{}

    collection := db.Collection("users")
    err := collection.FindOne(context.TODO(), bson.D{{"email", args.Email}}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    ctx.Value("log").(*logging.Logger).Debugf("GQL: Retrieved user by user [%s] : %v", *currentUserEmail, user)
    return &UserResolver{&user}, nil
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
        result = append(result, &UserResolver{&user})
    }

    if err := cur.Err(); err != nil {
      return nil, err
    }

    /*
    for _, item := range result {
        ctx.Value("log").(*logging.Logger).Debugf("GQL: User iteration result item:  %s", *item.u)
    }
    */

    return result, nil
}
