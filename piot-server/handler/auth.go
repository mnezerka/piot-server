package handler

import (
    "context"
    "encoding/json"
    //"fmt"
    "piot-server/config"
    "piot-server/config/db"
    "piot-server/model"
    "io/ioutil"
    //"log"
    "net/http"
    //jwt "github.com/dgrijalva/jwt-go"
    "github.com/mongodb/mongo-go-driver/bson"
    "golang.org/x/crypto/bcrypt"
)

func RegisterHandler(a *config.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

    w.Header().Set("Content-Type", "application/json")

    var user model.User

    // decode json from request body
    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(body, &user)

    var response model.ResponseResult

    if err != nil {
        response.Error = err.Error()
        json.NewEncoder(w).Encode(response)
        return 500, err
    }

    collection, err := db.GetDBCollection()

    if err != nil {
        response.Error = err.Error()
        json.NewEncoder(w).Encode(response)
        return 500, err
    }
    var result model.User
    err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

    if err != nil {
        if err.Error() == "mongo: no documents in result" {
            hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

            if err != nil {
                response.Error = "Error While Hashing Password, Try Again"
                json.NewEncoder(w).Encode(response)
                return 500, err
            }
            user.Password = string(hash)

            _, err = collection.InsertOne(context.TODO(), user)
            if err != nil {
                response.Error = "Error While Creating User, Try Again"
                json.NewEncoder(w).Encode(response)
                return 500, err
            }
            response.Result = "Registration Successful"
            json.NewEncoder(w).Encode(response)
            return 200, nil
        }

        response.Error = err.Error()
        json.NewEncoder(w).Encode(response)
        return 500, err
    }

    response.Result = "Username already Exists!!"
    json.NewEncoder(w).Encode(response)
    return 200, nil
}
