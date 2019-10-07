package handler

import (
    "fmt"
    "net/http"
    "mosquitto-auth/config"
)

var landingPage = []byte(fmt.Sprintf(`<html>
<head><title>PIOT Mosquitto Auth Server</title></head>
<body>
<h1>PIOT Mosquitto Auth Server</h1>
<p>Version: %s</p>
</body>
</html>
`, config.VersionString()))

func RootHandler(w http.ResponseWriter, r *http.Request) {

    // check http method, POST is required
    if r.Method != http.MethodGet {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write(landingPage)
}
