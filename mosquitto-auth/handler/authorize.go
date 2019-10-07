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
)

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

    // first, try to check static users
    switch packet.Username {
    case "test":


        if utils.GetMqttRootTopic(packet.Topic) == "test" {
            ctx.Value("log").(*logging.Logger).Debugf("Authorization passed for static user <%s> and topic <%s>", packet.Username, packet.Topic)
            return
        }
        ctx.Value("log").(*logging.Logger).Errorf("Authorization rejected for static user <%s> and topic <%s>", packet.Username, packet.Topic)
        http.Error(w, fmt.Sprintf("Authorization rejected for static user <%s> and topic <%s>", packet.Username, packet.Topic), 401)
        return
    case "mon":
        // TODO check also Acc attribute to allow only read
        if utils.GetMqttRootTopic(packet.Topic) == "$SYS" {
            ctx.Value("log").(*logging.Logger).Debugf("Authorization passed for static user <%s> and topic <%s>", packet.Username, packet.Topic)
            return
        }
        ctx.Value("log").(*logging.Logger).Errorf("Authorization rejected for static user <%s> and topic <%s>", packet.Username, packet.Topic)
        http.Error(w, fmt.Sprintf("Authorization rejected for static user <%s> and topic <%s>", packet.Username, packet.Topic), 401)
        return
    case "piot":
        if utils.GetMqttRootTopic(packet.Topic) == "org" {
            ctx.Value("log").(*logging.Logger).Debugf("Authorization passed for static user <%s> and topic <%s>", packet.Username, packet.Topic)
            return
        }
        ctx.Value("log").(*logging.Logger).Errorf("Authorization rejected for static user <%s> and topic <%s>", packet.Username, packet.Topic)
        http.Error(w, fmt.Sprintf("Authorization rejected for static user <%s> and topic <%s>", packet.Username, packet.Topic), 401)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Find user by id <%s>", packet.Username)

    // try to find user in database
    db := ctx.Value("db").(*mongo.Database)

    var user User

    err := db.Collection("users").FindOne(ctx, bson.M{"email": packet.Username}).Decode(&user)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Users service error: %v", err)
        http.Error(w, fmt.Sprintf("User identified as <%s> does not exist or provided credentials are wrong.", packet.Username), 401)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("User identified as <%s> was found.", packet.Username)

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
