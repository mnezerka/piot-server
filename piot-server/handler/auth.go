package handler

import (
    "context"
    "encoding/json"
    "errors"
    //"fmt"
    "regexp"
    "time"
    //"fmt"
    //"piot-server/config"
    //"piot-server/config/db"
    "piot-server/model"
    "github.com/op/go-logging"
    //"log"
    "net/http"
    jwt "github.com/dgrijalva/jwt-go"
    "github.com/mongodb/mongo-go-driver/bson"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/mongo"
)

// Create the JWT key used to create the signature
var JWT_KEY = []byte("my_secret_key")

const TOKEN_EXPIRATION = 5


func GetPasswordHash(pwd string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 5)
    return string(hash), err
}

func validateEmail(email string) bool {
    Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return Re.MatchString(email)
}

func WriteErrorResponse(w http.ResponseWriter, err error, status int) {
    var response model.ResponseResult
    response.Error = err.Error()
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}

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
        hash, err := GetPasswordHash(credentials.Password)
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

func LoginHandler() http.Handler {
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

        // try to find user in database
        var user model.User
        collection := db.Collection("users")
        err = collection.FindOne(ctx, bson.D{{"email", credentials.Email}}).Decode(&user)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf(err.Error())
            WriteErrorResponse(w, errors.New("User identified by this email does not exist or provided credentials are wrong!"), 401)
            return
        }

        // check if password is correct
        err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf(err.Error())
            WriteErrorResponse(w, errors.New("User identified by this email does not exist or provided credentials are wrong!"), 401)
            return
        }

        // Declare the expiration time of the token
        // here, we have kept it as 5 hours
        expirationTime := time.Now().Add(TOKEN_EXPIRATION * time.Hour)
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

        tokenString, err := token.SignedString(JWT_KEY)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf(err.Error())
            WriteErrorResponse(w, errors.New("Error while generating token, try again"), 500)
            return
        }

        var response model.Token
        response.Token = tokenString

        ctx.Value("log").(*logging.Logger).Debugf("Successfull login: %s", user.Email)

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)

        return
    })
}

/*
func RefreshHandler(w http.ResponseWriter, r *http.Request) {

    // try to parse JWT from Authorization header
    tokenString := r.Header.Get("Authorization")

    claims := &Claims{}

    // Parse the JWT string and store the result in `claims`.
    // Note that we are passing the key in this method as well. This method will return an error
    // if the token is invalid (if it has expired according to the expiry time we set on sign in),
    // or if the signature does not match
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("Unexpected signing method")
        }
        return jwtKey, nil
    })

    if !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // We ensure that a new token is not issued until enough time has elapsed
    // In this case, a new token will only be issued if the old token is within
    // 1 hour of expiry. Otherwise, return a bad request status
    if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 1 * time.Hour {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // Now, create a new token for the current use, with a renewed expiration time
    expirationTime := time.Now().Add(TOKEN_EXPIRATION * time.Hour)
    claims.ExpiresAt = expirationTime.Unix()
    newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    newTokenString, err := newToken.SignedString(jwtKey)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    var response model.Token
    response.Token = newTokenString

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

}
*/
