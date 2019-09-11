/*
 Resources:
    https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/
*/

package main

import (
    "log"
    "net/http"
    "os"
    "golang.org/x/net/context"
    "github.com/urfave/cli"
    //"github.com/gorilla/mux"
    //"github.com/gorilla/handlers"
    "piot-server/handler"
    "piot-server/config"
    "piot-server/service"
    "piot-server/db"
)

const LOG_FORMAT = "%{color}%{time:2006/01/02 15:04:05 -07:00 MST} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"

func runServer(c *cli.Context) {

    // create global context for all handlers
    db, err := db.GetDB(c.GlobalString("mongodb-uri"))
    fatalfOnError(err, "Failed to open database on %s", c.GlobalString("mongodb-uri"))

    ctx := context.Background()
    log := service.NewLogger(LOG_FORMAT, false)

    log.Infof("Starting PIOT server %s", config.VersionString())

    //ctx := &config.AppContext{Db: db}

    ctx = context.WithValue(ctx, "db", db)
    ctx = context.WithValue(ctx, "log", log)

    http.HandleFunc("/", handler.RootHandler)

    // endpoint for registration of new user
    http.Handle("/register", handler.AddContext(ctx, handler.Logging(handler.Registration())))

    // endpoint for authentication - token is generaged
    http.Handle("/login", handler.AddContext(ctx, handler.Logging(handler.LoginHandler())))

    //r := mux.NewRouter()

    //r.HandleFunc("/", handler.RootHandler)

    //r.PathPrefix("/register").Handler(handler.AppHandler{ctx, handler.RegisterHandler}).
    //    Methods("POST")

    //r.PathPrefix("/signin").Handler(handler.AppHandler{ctx, handler.SigninHandler}).
    //    Methods("POST")

    //r.HandleFunc("/refresh", handler.RefreshHandler)

    log.Infof("Listening on %s...", c.GlobalString("bind-address"))
    //err = http.ListenAndServe(c.GlobalString("bind-address"), handlers.LoggingHandler(os.Stdout, r))
    err = http.ListenAndServe(c.GlobalString("bind-address"), nil)
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

    app.Name = "PIOT Server"
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
            Name:   "mqtt-uri,q",
            Usage:  "Endpoint for the Mosquitto message broker",
            EnvVar: "MQTT_URI",
            Value:  "tcp://localhost:1883",
        },
        cli.StringFlag{
            Name:   "bind-address,b",
            Usage:  "Listen address for API HTTP endpoint",
            Value:  "0.0.0.0:9096",
            EnvVar: "BIND_ADDRESS",
        },
        cli.StringFlag{
            Name:   "mongodb-uri,m",
            Usage:  "URI for the mongo database",
            Value:  "mongodb://localhost:27017",
            EnvVar: "MONGODB_URI",
        },
    }

    app.Run(os.Args)
}
