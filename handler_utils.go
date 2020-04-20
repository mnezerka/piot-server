package main

import (
    "encoding/json"
    "piot-server/model"
    "net/http"
)

func WriteErrorResponse(w http.ResponseWriter, err error, status int) {
    var response model.ResponseResult
    response.Error = err.Error()
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}
