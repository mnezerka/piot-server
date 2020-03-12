package handler

import (
    "context"
    "encoding/json"
    "errors"
    "time"
    "github.com/mnezerka/go-piot/config"
    "github.com/mnezerka/go-piot/model"
    "github.com/op/go-logging"
    "net/http"
    jwt "github.com/dgrijalva/jwt-go"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "golang.org/x/crypto/bcrypt"
)

type Login struct {
    log *logging.Logger
    db *mongo.Database
    params *config.Parameters
}

func NewLogin(log *logging.Logger, db *mongo.Database, params *config.Parameters) *Login {
    return &Login{log: log, db: db, params: params}
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

    // try to find user in database
    var user model.User
    collection := l.db.Collection("users")
    err = collection.FindOne(context.TODO(), bson.D{{"email", credentials.Email}}).Decode(&user)
    if err != nil {
        l.log.Errorf(err.Error())
        WriteErrorResponse(w, errors.New("User identified by this email does not exist or provided credentials are wrong!"), 401)
        return
    }

    // check if password is correct
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
    if err != nil {
        l.log.Errorf(err.Error())
        WriteErrorResponse(w, errors.New("User identified by this email does not exist or provided credentials are wrong!"), 401)
        return
    }

    // Declare the expiration time of the token
    // here, we have kept it as 5 hours
    expirationTime := time.Now().Add(l.params.JwtTokenExpiration)
    l.log.Debugf("Setting expiration to %v (%d)", expirationTime, expirationTime.Unix())

    // Create the JWT claims, which includes the username and expiry time
    claims := &model.Claims{
        Email: user.Email,
        StandardClaims: jwt.StandardClaims{
            // In JWT, the expiry time is expressed as unix milliseconds
            ExpiresAt: expirationTime.Unix(),
        },
    }

    // generate new jwt token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    //ctx.Value("log").(*logging.Logger).Debugf("JWT Pass: %s", params.JwtPassword)

    tokenString, err := token.SignedString([]byte(l.params.JwtPassword))
    if err != nil {
        l.log.Errorf(err.Error())
        WriteErrorResponse(w, errors.New("Error while encrypting token, try again"), 500)
        return
    }

    //ctx.Value("log").(*logging.Logger).Debugf("JWT Token: %s", tokenString)

    var response model.Token
    response.Token = tokenString

    l.log.Debugf("Successfull login: %s", user.Email)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

    return
}
