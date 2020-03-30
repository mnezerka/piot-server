package main

import (
    "fmt"
    "errors"
    "net/http"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "strings"
    jwt "github.com/dgrijalva/jwt-go"
    "piot-server/model"
    "piot-server/config"
)

type AuthHandler struct {
    log *logging.Logger
    cfg *config.Parameters
    users *Users
    handler http.Handler
}

func NewAuthHandler(log *logging.Logger, cfg *config.Parameters, users *Users, handler http.Handler) *AuthHandler {
    return &AuthHandler{log: log, cfg: cfg, users: users, handler: handler}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    isAuthorized := false

    var tokenString string

    ctx := r.Context()

    // 1. first - try to get auth token from query parameter "token"
    keys, ok := r.URL.Query()["token"]
    if !ok || len(keys) < 1 {
        // second - try to get auth token from authorization header
        auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
        if len(auth) != 2 || auth[0] != "Bearer" {
            WriteErrorResponse(w, errors.New("Invalid or missing authorization header"), 401)
            return
        }
        tokenString = auth[1]
    } else {
        tokenString = keys[0]
    }

    // 2. we have token string, let's validate it

    // Initialize a new instance of `Claims`
    claims := &model.Claims{}

    // Parse the JWT string and store the result in `claims`.
    // Note that we are passing the key in this method as well. This method will return an error
    // if the token is invalid (if it has expired according to the expiry time we set on sign in),
    // or if the signature does not match
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(h.cfg.JwtPassword), nil
    })

    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            WriteErrorResponse(w, err, http.StatusUnauthorized)
            return
        }
        WriteErrorResponse(w, err, http.StatusUnauthorized)
        return
    }

    h.log.Debugf("Authentication passed")

    if !token.Valid {
        WriteErrorResponse(w, errors.New("Token is not valid"), 401)
        return
    }

    h.log.Debugf("Token is valid, email: %s", claims.Email)


    // 3. Find user in database to prepare user profile

    user, err := h.users.FindByEmail(claims.Email)

    ctx = context.WithValue(ctx, "user_email", &claims.Email)
    ctx = context.WithValue(ctx, "is_authorized", isAuthorized)

    ctx = context.WithValue(ctx, "profile", &model.UserProfile{
        claims.Email,       // email
        false,              // is admin
        nil,
        user.Orgs,
    })

    h.handler.ServeHTTP(w, r.WithContext(ctx))
}
