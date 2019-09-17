package main

import (
    "log"
    "net/http"
    "os"
    "golang.org/x/net/context"
    "github.com/urfave/cli"
    "piot-mqtt-adapter/handler"
    "piot-mqtt-adapter/config"
    "piot-mqtt-adapter/service"
)

const LOG_FORMAT = "%{color}%{time:2006/01/02 15:04:05 -07:00 MST} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"

func runServer(c *cli.Context) {

    // create global context for all handlers
    ctx := context.Background()

    // create global logger for all handlers
    log := service.NewLogger(LOG_FORMAT, true)

    log.Infof("Starting PIOT MQTT adapter %s", config.VersionString())

    ctx = context.WithValue(ctx, "log", log)

    //http.HandleFunc("/", handler.RootHandler)
    http.Handle("/", handler.AddContext(ctx, handler.Logging(handler.RootHandler())))

    log.Infof("Listening on %s...", c.GlobalString("bind-address"))
    err := http.ListenAndServe(c.GlobalString("bind-address"), nil)
    fatalfOnError(err, "Failed to bind on %s: ", c.GlobalString("bind-address"))
}

func fatalfOnError(err error, msg string, args ...interface{}) {
    if err != nil {
        log.Fatalf(msg, args...)
        os.Exit(1)
    }
}

func main() {
    app := cli.NewApp()

    app.Name = "PIOT MQTT Adapter"
    app.Version = config.VersionString()
    app.Authors = []cli.Author{
        {
            Name:  "Michal Nezerka",
            Email: "michal.nezerka@gmail.com",
        },
    }
    app.Usage = "Management of Pavoucek IOT things"
    app.Action = runServer
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:   "piot-server-uri,p",
            Usage:  "Endpoint for the Mosquitto message broker",
            EnvVar: "PIOT_SERVER_URI",
            Value:  "tcp://localhost:9096",
        },
        cli.StringFlag{
            Name:   "bind-address,b",
            Usage:  "Listen address for API HTTP endpoint",
            Value:  "0.0.0.0:9097",
            EnvVar: "BIND_ADDRESS",
        },
        cli.StringFlag{
            Name:   "mqtt-uri,q",
            Usage:  "Endpoint for the Mosquitto message broker",
            EnvVar: "MQTT_URI",
            Value:  "tcp://localhost:1883",
        },
    }

    app.Run(os.Args)
}
