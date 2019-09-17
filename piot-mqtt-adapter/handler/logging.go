package handler

import (
    "github.com/op/go-logging"
    "net/http"
)

func Logging(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        log := ctx.Value("log").(*logging.Logger)
        log.Infof("%s %s %s %s %s", r.RemoteAddr, r.Method, r.URL, r.Proto, r.UserAgent())
        h.ServeHTTP(w, r)
    })
}
