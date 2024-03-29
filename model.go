package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Used to read the email and password from the token
// request body (signin use case)
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Used to serialize token as a response to authentication request
type Token struct {
	Token string `json:"token"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

// Struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields
// like expiry time
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// this becomes part of context that is propagated to all handlers - e.g. graphql
type UserProfile struct {
	Id      primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Email   string               `json:"email" bson:"email"`
	IsAdmin bool                 `json:"is_admin" bson:"is_admin"`
	OrgId   primitive.ObjectID   `json:"org_id" bson:"org_id"`
	OrgIds  []primitive.ObjectID `json:"orgs" bson:"orgs"`
}
