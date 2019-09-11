package service

import (
    "context"
    "errors"
    "github.com/op/go-logging"
    "piot-server/model"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"

)

type UserService struct {
    db          *mongo.Database
    log         *logging.Logger
}

func NewUserService(db *mongo.Database, log *logging.Logger) *UserService {
    return &UserService{db: db, log: log}
}

func (u *UserService) FindByEmail(email string) (*model.User, error) {
    user := &model.User{}

    collection := u.db.Collection("users")
    err := collection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
    if err != nil {
        return nil, errors.New("User not found")
    }

    return user, nil
}

func (u *UserService) List() ([]*model.User, error) {
    users := make([]*model.User, 0)

    collection := u.db.Collection("users")
    cur, err := collection.Find(context.TODO(), bson.D{})
    if err != nil { return nil, errors.New("Error while fetching users from database")}
    defer cur.Close(context.TODO())
    for cur.Next(context.TODO()) {
        user := &model.User{}
        err := cur.Decode(&user)
        if err != nil { return nil, errors.New("Error while fetching users from database")}

    }

    if err := cur.Err(); err != nil {
        return nil, err
    }

    return users, nil
}

/*
func (u *UserService) Count() (int, error) {
    var count int
    userSQL := `SELECT count(*) FROM users`
    err := u.db.Get(&count, userSQL)
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (u *UserService) ComparePassword(userCredentials *model.UserCredentials) (*model.User, error) {
    user, err := u.FindByEmail(userCredentials.Email)
    if err != nil {
        return nil, errors.New(configuration.UnauthorizedAccess)
    }
    if result := user.ComparePassword(userCredentials.Password); !result {
        return nil, errors.New(configuration.UnauthorizedAccess)
    }
    return user, nil
}
*/
