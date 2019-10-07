package handler

import (
    //"bytes"
    //"net/http/httputil"
    "github.com/op/go-logging"
    //"io/ioutil"
    "net/http"
)

func Logging(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        log := ctx.Value("log").(*logging.Logger)
        log.Infof("%s %s %s %s %s", r.RemoteAddr, r.Method, r.URL, r.Proto, r.UserAgent())
        //requestDump, _ := httputil.DumpRequest(r, true)
        //ctx.Value("log").(*logging.Logger).Debugf("Incoming mosquitto auth request:\n%s", requestDump)


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
        h.ServeHTTP(w, r)
    })
}
