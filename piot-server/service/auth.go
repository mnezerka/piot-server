package service

import (
    //"fmt"
    //jwt "github.com/dgrijalva/jwt-go"
    "github.com/op/go-logging"
    "time"
)

type AuthService struct {
    signedSecret        *string
    expiredTimeInSecond *time.Duration
    log                 *logging.Logger
}

/*
func NewAuthService(config *configuration.Config, log *logging.Logger) *AuthService {
    return &AuthService{&config.JWTSecret, &config.JWTExpireIn, log}
}

func (a *AuthService) SignJWT(user *model.User) (*string, error) {

    // Declare the expiration time of the token
    // here, we read this value (in seconds) from configuration
    expirationTime := time.Now().Add(time.Second * *a.expiredTimeInSecond)

    // Create the JWT claims, which includes the user id and expiry time
    claims := &model.Claims{
        Id: user.ID,
        StandardClaims: jwt.StandardClaims{
            // In JWT, the expiry time is expressed as unix milliseconds
            ExpiresAt: expirationTime.Unix(),
        },
    }

    // Declare the token with the algorithm used for signing, and the claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString([]byte(*a.signedSecret))

    return &tokenString, err
}

// validate token provided as string, return parsed token
func (a *AuthService) ValidateJWT(tokenString *string) (*jwt.Token, *model.Claims, error) {

    // Initialize a new instance of `Claims`
    claims := &model.Claims{}

    //token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
    token, err := jwt.ParseWithClaims(*tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("    unexpected signing method: %v", token.Header["alg"])
        }

        return []byte(*a.signedSecret), nil
    })

    return token, claims, err
}
*/
