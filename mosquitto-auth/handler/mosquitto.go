package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/op/go-logging"
    "mosquitto-auth/utils"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"
)

type MosquittoAuthUser struct {
    Username    string `json:"username"`
    Password    string `json:"password"`
}

type MosquittoAuthAcl struct {
    Acc         int `json:"acc"`
    ClientId    string `json:"clientid"`
    Topic       string `json:"topic"`
    Username    string `json:"username"`
}

// Represents org as stored in database
type Org struct {
    Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string `json:"name"`
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

    // try to find user in database
    db := ctx.Value("db").(*mongo.Database)

    var user User
    collection := db.Collection("users")
    err := collection.FindOne(ctx, bson.D{{"email", packet.Username}}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf(err.Error())
        http.Error(w, fmt.Sprintf("User identified by email %s does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("User %s exists", packet.Username)

    // check if password is correct
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(packet.Password))
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf(err.Error())
        http.Error(w, fmt.Sprintf("User identified by email %s does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Authentication for user %s passed", packet.Username)
}

type AuthenticateSuperUser struct { }

func (h *AuthenticateSuperUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ctx.Value("log").(*logging.Logger).Debugf("Request for superuser authentication, denying (not supported)")
    http.Error(w, "Superuser role Not supported in PIOT", 401)
}


type Authorize struct { }

func (h *Authorize) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    // check http method, POST is required
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // try to decode packet
    var packet MosquittoAuthAcl
    if err := json.NewDecoder(r.Body).Decode(&packet); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx := r.Context()
    ctx.Value("log").(*logging.Logger).Debugf("Acl request for user %s, topic: %s, client: %s, access type: %d", packet.Username, packet.Topic, packet.ClientId, packet.Acc)

    ctx.Value("log").(*logging.Logger).Debugf("Get user by email: %s", packet.Username)

    // try to find user in database
    db := ctx.Value("db").(*mongo.Database)

    var user User

    err := db.Collection("users").FindOne(ctx, bson.M{"email": packet.Username}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Users service error: %v", err)
        http.Error(w, fmt.Sprintf("User identified by email %s does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("User identified by email <%s> was found.", packet.Username)

    // fetch user orgs

    // filter orgusers to current (single) user
    stage_match := bson.M{"$match": bson.M{"user_id": user.Id}}

    // find orgs details
    stage_lookup := bson.M{"$lookup": bson.M{"from": "orgs", "localField": "org_id", "foreignField": "_id", "as": "orgs"}}

    // expand orgs
    stage_unwind := bson.M{"$unwind": "$orgs"}

    // replace root
    stage_new_root := bson.M{"$replaceWith": "$orgs"}

    pipeline := []bson.M{stage_match, stage_lookup, stage_unwind, stage_new_root}

    cur, err := db.Collection("orgusers").Aggregate(ctx, pipeline)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Error while querying user orgs: %v", err)
        http.Error(w, "Fetching of user orgs failed", 500)
        return
    }
    defer cur.Close(ctx)

    for cur.Next(ctx) {
        var org Org
        if err := cur.Decode(&org); err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("Error while querying user orgs: %v", err)
            http.Error(w, "Fetching of user orgs failed", 500)
            return
        }
        user.Orgs = append(user.Orgs, org)
    }

    if err := cur.Err(); err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Error while querying user orgs: %v", err)
        http.Error(w, "Fetching of user orgs failed", 500)
        return
    }

    // extract org from topic name and check if user is member of given org
    orgName := utils.GetMqttTopicOrg(packet.Topic)
    if orgName != "" {
        for _, userOrg := range user.Orgs {
            if userOrg.Name == orgName {
                ctx.Value("log").(*logging.Logger).Debugf("Topic is matching user org (%s) -> authorization passed", orgName)
                return
            }
        }
    }

    ctx.Value("log").(*logging.Logger).Debugf("No org matching topic %s -> authorization failed", orgName)
    http.Error(w, fmt.Sprintf("User is not assigned to organization %s", orgName), 401)
}
