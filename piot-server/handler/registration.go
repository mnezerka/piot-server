package handler

import (
    "context"
    "encoding/json"
    "errors"
    "time"
    "piot-server/model"
    "piot-server/utils"
    "github.com/op/go-logging"
    "net/http"
    "github.com/mongodb/mongo-go-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func Registration() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        ctx := r.Context()
        db := ctx.Value("db").(*mongo.Database)

        // check http method, POST is required
        if r.Method != http.MethodPost {
            WriteErrorResponse(w, errors.New("Only POST method is allowed"), http.StatusMethodNotAllowed)
            return
        }

        // decode json from request body
        var credentials model.Credentials
        err := json.NewDecoder(r.Body).Decode(&credentials)
        if err != nil {
            WriteErrorResponse(w, err, 400)
            return
        }

        // check required attributes
        if len(credentials.Email) == 0 {
            WriteErrorResponse(w, errors.New("Email field is empty or not specified!"), 400)
            return
        }
        if len(credentials.Password) == 0 {
            WriteErrorResponse(w, errors.New("Password field is empty or not specified!"), 400)
            return
        }
        if !validateEmail(credentials.Email) {
            WriteErrorResponse(w, errors.New("Email field has wrong format!"), 400)
            return
        }

        // try to find existing user
        var user model.User
        collection := db.Collection("users")
        err = collection.FindOne(context.TODO(), bson.D{{"email", credentials.Email}}).Decode(&user)
        if err == nil {
            WriteErrorResponse(w, errors.New("User identified by this email already exists!"), 409)
            return
        }

        // generate hash for given password (we don't store passwords in plain form)
        hash, err := utils.GetPasswordHash(credentials.Password)
        if err != nil {
            WriteErrorResponse(w, errors.New("Error while hashing password, try again"), 500)
            return
        }
        user.Email = credentials.Email
        user.Password = hash
        user.Created = int32(time.Now().Unix())

        // user does not exist -> create new one
        _, err = collection.InsertOne(context.TODO(), user)
        if err != nil {
            WriteErrorResponse(w, errors.New("User while creating user, try again"), 500)
            return
        }

        var response model.ResponseResult
        response.Result = "Registration successful"

        ctx.Value("log").(*logging.Logger).Debugf("User is registered: %s", user.Email)

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)

        return
    })
}
