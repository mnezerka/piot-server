package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "piot-mqtt-adapter/config"
    "piot-mqtt-adapter/model"
)

var landingPage = []byte(fmt.Sprintf(`<html>
<head><title>PIOT MQTT Adapter</title></head>
<body>
<h1>PIOT MQTT Adapter</h1>
<p>Version: %s</p>
</body>
</html>
`, config.VersionString()))

func RootHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            w.Write(landingPage)
        case "POST":
            HandlePost(w, r)
        default:
            fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
        }
    })
}

func HandlePost(w http.ResponseWriter, r *http.Request) {

    var request model.Request

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    //ctx := r.Context()
    if len(request.Device) == 0 {
        http.Error(w, "Missing device ID", http.StatusBadRequest)
        return
    }

    fmt.Fprintf(w, "Post from website! %v", request)
}
