package handler

import (
    "encoding/json"

    //"log"
    "net/http"
    "piot-server/config"
    "piot-server/model"
)

// Application handler
type AppHandler struct {
    Context *config.AppContext
    HandleFunc func(*config.AppContext, http.ResponseWriter, *http.Request) (int, error)
}

// Our ServeHTTP method is mostly the same, and also has the ability to
// access our *appContext's fields (templates, loggers, etc.) as well.
func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Updated to pass ah.appContext as a parameter to our handler type.
    //status, err := ah.HandleFunc(ah.Context, w, r)
    status, err := ah.HandleFunc(ah.Context, w, r)
    //log.Printf("HTTP %d: %q", status, err)

    if err != nil {
        var response model.ResponseResult
        response.Error = err.Error()
        http.Error(w, http.StatusText(status), status)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

