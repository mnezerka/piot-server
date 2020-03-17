package handler

import (
    "net/http"
)

type CORSHandler struct {
    handler http.Handler
}

func NewCORSHandler(handler http.Handler) *CORSHandler {
    return &CORSHandler{handler: handler}
}

func (h *CORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    // Allow from any origin
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Max-Age", "86400")    // cache for 1 day

    // Handle pre-flight OPTIONS requests
    if r.Method == http.MethodOptions {

        if r.Header.Get("Access-Control-Request-Method") != "" {
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
        }

        if r.Header.Get("Access-Control-Request-Headers") != "" {
            w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
        }

        return
    }

    h.handler.ServeHTTP(w, r)
}
