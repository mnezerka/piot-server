package handler

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    //"net/http/httputil"
    "github.com/op/go-logging"
    "piot-server/service"
    "piot-server/utils"
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

type MosquittoAuth struct { }

func (h *MosquittoAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    ctx := r.Context()
    ctx.Value("log").(*logging.Logger).Debugf("Incoming mosquitto auth request path:%s", r.URL.Path)
    //requestDump, _ := httputil.DumpRequest(r, true)
    //ctx.Value("log").(*logging.Logger).Debugf("Incoming mosquitto auth request:\n%s", requestDump)

    // check http method, POST is required
    if r.Method != http.MethodPost {
        WriteErrorResponse(w, errors.New("Only POST method is allowed"), http.StatusMethodNotAllowed)
        return
    }

    switch r.URL.Path {
        case "/mosquitto-auth-user":

            // try to decode packet
            var packet MosquittoAuthUser
            if err := json.NewDecoder(r.Body).Decode(&packet); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            ctx.Value("log").(*logging.Logger).Debugf("Authenticating user %s", packet.Username)

            auth := ctx.Value("auth").(*service.Auth)

            if err := auth.AuthUser(ctx, packet.Username, packet.Password); err != nil {
                http.Error(w, err.Error(), 401)
            }
        case "/mosquitto-auth-superuser":

            ctx.Value("log").(*logging.Logger).Debugf("Request for superuser authentication, denying (not supported)")
            http.Error(w, "Superuser role Not supported in PIOT", 401)

        case "/mosquitto-auth-acl":

            // try to decode packet
            var packet MosquittoAuthAcl
            if err := json.NewDecoder(r.Body).Decode(&packet); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            ctx.Value("log").(*logging.Logger).Debugf("Acl request for user %s, topic: %s, client: %s, access type: %d", packet.Username, packet.Topic, packet.ClientId, packet.Acc)

            users := ctx.Value("users").(*service.Users)

            user, err := users.FindByEmail(ctx, packet.Username)
            if err != nil {
                http.Error(w, err.Error(), 401)
                return
            }

            ctx.Value("log").(*logging.Logger).Debugf("Fetched user: %s (%d orgs)", user.Email, len(user.Orgs))

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

        default:
            ctx.Value("log").(*logging.Logger).Errorf("Unkown path for mosquitto authentication: %s", r.URL.Path)
            http.Error(w, "Unknown path", 403)
    }

}
