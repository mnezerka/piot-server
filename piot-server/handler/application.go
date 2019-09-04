package handler

import (
    "log"
    "net/http"
    "piot-server/config"
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
    status, err := ah.HandleFunc(ah.Context, w, r)
    if err != nil {
        log.Printf("HTTP %d: %q", status, err)
        switch status {
        case http.StatusNotFound:
            http.NotFound(w, r)
            // And if we wanted a friendlier error page, we can
            // now leverage our context instance - e.g.
            // err := ah.renderTemplate(w, "http_404.tmpl", nil)
        case http.StatusInternalServerError:
            http.Error(w, http.StatusText(status), status)
        default:
            http.Error(w, http.StatusText(status), status)
        }
    }
}

