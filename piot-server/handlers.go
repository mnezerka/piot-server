package main

import (
	"fmt"
	"net/http"
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
`, versionString()))

func serveVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(landingPage)
}
