package model


// Represents user as stored in database
type User struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
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
