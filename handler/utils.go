package handler

import (
    "encoding/json"
    "regexp"
    "piot-server/model"
    "net/http"
)

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
