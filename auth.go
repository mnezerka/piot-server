package main

import (
    "context"
    "fmt"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "golang.org/x/crypto/bcrypt"
    "piot-server/model"
)

type Auth struct { }

func (a *Auth) AuthUser(ctx context.Context, email, password string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Authenticate user: %s", email)

    db := ctx.Value("db").(*mongo.Database)

    // try to find user in database
    var user model.User
    collection := db.Collection("users")
    err := collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf(err.Error())
        return fmt.Errorf("User identified by email %s does not exist or provided credentials are wrong.", email)
    }

    ctx.Value("log").(*logging.Logger).Debugf("User %s exists", email)

    // check if password is correct
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf(err.Error())
        return fmt.Errorf("User identified by email %s does not exist or provided credentials are wrong.", email)
    }

    ctx.Value("log").(*logging.Logger).Debugf("Authentication for user %s passed", email)

    return nil
}
