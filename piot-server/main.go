package main

import (
    "log"
    "net/http"
    "os"
    "github.com/urfave/cli"
    //mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
    appName = "PIOT Server"
)

var (
    topics = map[string]string{
        "test":            "The total number of bytes received since the broker started.",
    }
)

func main() {
    app := cli.NewApp()

    app.Name = appName
    app.Version = versionString()
    app.Authors = []cli.Author{
        {
            Name:  "Michal Nezerka",
            Email: "michal.nezerka@gmail.com",
        },
    }
    app.Usage = "Management of IOT things"
    app.Action = runServer
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:   "endpoint,e",
            Usage:  "Endpoint for the Mosquitto message broker",
            EnvVar: "BROKER_ENDPOINT",
            Value:  "tcp://127.0.0.1:1883",
        },
        cli.StringFlag{
            Name:   "bind-address,b",
            Usage:  "Listen address for API HTTP endpoint",
            Value:  "0.0.0.0:9096",
            EnvVar: "BIND_ADDRESS",
        },
    }

    app.Run(os.Args)
}

func runServer(c *cli.Context) {
    log.Printf("Starting PIOT server %s", versionString())

    http.HandleFunc("/", serveVersion)

    log.Printf("Listening on %s...", c.GlobalString("bind-address"))
    err := http.ListenAndServe(c.GlobalString("bind-address"), nil)
    fatalfOnError(err, "Failed to bind on %s: ", c.GlobalString("bind-address"))
}

func fatalfOnError(err error, msg string, args ...interface{}) {
    if err != nil {
        log.Fatalf(msg, args...)
        os.Exit(1)
    }
}
