package main

import (
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, err error, status int) {
	var response ResponseResult
	response.Error = err.Error()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
