package handler

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
