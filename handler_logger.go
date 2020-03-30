package main

import (
    //"bytes"
    "github.com/op/go-logging"
    //"io/ioutil"
    "net/http"
)

type LoggingHandler struct {
    log *logging.Logger
    handler http.Handler
}

func NewLoggingHandler(log *logging.Logger, handler http.Handler) *LoggingHandler {
    return &LoggingHandler{log: log, handler: handler}
}

func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.log.Debugf("%s %s %s %s %s", r.RemoteAddr, r.Method, r.URL, r.Proto, r.UserAgent())
    /*
    TODO - get info of debug mode directly from logger
    if l.DebugMode {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Errorf("Reading request body error: %s", err)
        }
        reqStr := ioutil.NopCloser(bytes.NewBuffer(body))
        log.Debugf("Request body : %v", reqStr)
        r.Body = reqStr
    }
    */
    h.handler.ServeHTTP(w, r)
}
