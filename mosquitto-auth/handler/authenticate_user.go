package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/op/go-logging"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"
)

type MosquittoAuthUser struct {
    Username    string `json:"username"`
    Password    string `json:"password"`
}

// Represents user as stored in database
type User struct {
    Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Email       string `json:"email"`
    Password    string `json:"password"`
    Orgs        []Org  `json:"orgs"`
}

type AuthenticateUser struct { }

func (h *AuthenticateUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    ctx := r.Context()

    // check http method, POST is required
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // try to decode packet
    var packet MosquittoAuthUser
    if err := json.NewDecoder(r.Body).Decode(&packet); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Authenticating user %s", packet.Username)

    // first - try static users
    switch packet.Username {
    case "test":
        if ctx.Value("test-pwd") != "" &&  ctx.Value("test-pwd") == packet.Password {
            ctx.Value("log").(*logging.Logger).Debugf("User <%s> authenticated as static", packet.Username)
            return
        }
        ctx.Value("log").(*logging.Logger).Errorf("Static user <%s> authenticated failed ", packet.Username)
        http.Error(w, fmt.Sprintf("User identified as <%s> does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    case "mon":
        if ctx.Value("mon-pwd") != "" &&  ctx.Value("mon-pwd") == packet.Password {
            ctx.Value("log").(*logging.Logger).Debugf("User <%s> authenticated as static", packet.Username)
            return
        }
        ctx.Value("log").(*logging.Logger).Errorf("Static user <%s> authenticated failed ", packet.Username)
        http.Error(w, fmt.Sprintf("User identified as <%s> does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    case "piot":
        if ctx.Value("piot-pwd") != "" &&  ctx.Value("piot-pwd") == packet.Password {
            ctx.Value("log").(*logging.Logger).Debugf("User <%s> authenticated as static", packet.Username)
            return
        }
        ctx.Value("log").(*logging.Logger).Errorf("Static user <%s> authenticated failed ", packet.Username)
        http.Error(w, fmt.Sprintf("User identified as <%s> does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    }


    // try to find user in database
    db := ctx.Value("db").(*mongo.Database)

    var user User
    collection := db.Collection("users")
    err := collection.FindOne(ctx, bson.D{{"email", packet.Username}}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf(err.Error())
        http.Error(w, fmt.Sprintf("User identified as <%s> does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("User %s exists", packet.Username)

    // check if password is correct
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(packet.Password))
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf(err.Error())
        http.Error(w, fmt.Sprintf("User identified as <%s> does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Authentication for user %s passed", packet.Username)
}
