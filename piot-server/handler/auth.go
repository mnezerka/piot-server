package handler

import (
    "context"
    "encoding/json"
    "fmt"
    "go-login/config/db"
    "go-login/model"
    "io/ioutil"
    "log"
    "net/http"

    jwt "github.com/dgrijalva/jwt-go"
    "github.com/mongodb/mongo-go-driver/bson"
    "golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")

    var user model.User

    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(body, &user)
    var res model.ResponseResult
    if err != nil {
        res.Error = err.Error()
        json.NewEncoder(w).Encode(res)
        return
    }

    collection, err := db.GetDBCollection()

    if err != nil {
        res.Error = err.Error()
        json.NewEncoder(w).Encode(res)
        return
    }
    var result model.User
    err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

    if err != nil {
        if err.Error() == "mongo: no documents in result" {
            hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

            if err != nil {
                res.Error = "Error While Hashing Password, Try Again"
                json.NewEncoder(w).Encode(res)
                return
            }
            user.Password = string(hash)

            _, err = collection.InsertOne(context.TODO(), user)
            if err != nil {
                res.Error = "Error While Creating User, Try Again"
                json.NewEncoder(w).Encode(res)
                return
            }
            res.Result = "Registration Successful"
            json.NewEncoder(w).Encode(res)
            return
        }

        res.Error = err.Error()
        json.NewEncoder(w).Encode(res)
        return
    }

    res.Result = "Username already Exists!!"
    json.NewEncoder(w).Encode(res)
    return
}
