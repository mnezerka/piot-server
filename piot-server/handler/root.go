package handler

import (
    "fmt"
    "net/http"
    "piot-server/config"
)

/*
 * Root and Healthcheck
 */

var landingPage = []byte(fmt.Sprintf(`<html>
<head><title>Mosquitto exporter</title></head>
<body>
<h1>PIOT Server</h1>
<p>Version: %s</p>
</body>
</html>
`, config.VersionString()))

func RootHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write(landingPage)
}
