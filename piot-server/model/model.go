package model

import (
    jwt "github.com/dgrijalva/jwt-go"
)


// Represents user as stored in database
type User struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
    Created   int32  `json:"created"`
}

type UserProfile struct {
    Email     string `json:"email"`
}

// Represents customer as stored in database
type Customer struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Created     int32  `json:"created"`
}

// Used to read the username and password from the request body
// for signin and authentication requests
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


