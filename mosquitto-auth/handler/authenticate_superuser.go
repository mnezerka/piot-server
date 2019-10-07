package handler

import (
    "net/http"
    "github.com/op/go-logging"
)

type AuthenticateSuperUser struct { }

func (h *AuthenticateSuperUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ctx.Value("log").(*logging.Logger).Debugf("Request for superuser authentication, denying (not supported)")
    http.Error(w, "Superuser role Not supported in PIOT", 401)
}
