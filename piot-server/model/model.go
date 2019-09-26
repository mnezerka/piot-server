package model

import (
    jwt "github.com/dgrijalva/jwt-go"
)

// Used to read the email and password from the token
// request body (signin use case)
type Credentials struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
}

// Used to serialize token as a response to authentication request
type Token struct {
    Token     string `json:"token"`
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


