package main

import (
    "log"
    "net/http"
    "os"
    "github.com/urfave/cli"
    "mosquitto-auth/handler"
    "mosquitto-auth/config"
    piotcontext "mosquitto-auth/context"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/op/go-logging"
)

func runServer(c *cli.Context) {

    // create global context for all handlers
    contextOptions := piotcontext.NewContextOptions()
    contextOptions.DbUri = c.GlobalString("mongodb-uri")
    contextOptions.DbName = "piot"
    contextOptions.LogLevel = c.GlobalString("log-level")
    contextOptions.TestPassword = c.GlobalString("test-password")
    contextOptions.MonPassword = c.GlobalString("mon-password")
    contextOptions.PiotPassword = c.GlobalString("piot-password")

    ctx := piotcontext.NewContext(contextOptions)

    // Auto disconnect from mongo
    defer ctx.Value("dbClient").(*mongo.Client).Disconnect(ctx)

    logger := ctx.Value("log").(*logging.Logger)
    logger.Infof("Starting PIOT mosquitto auth server %s", config.VersionString())

    /////////////// HANDLERS

    http.HandleFunc("/", handler.RootHandler)

    // endpoints for mosquitto authentication requests
    http.Handle("/mosquitto-auth-user", handler.CORS(handler.AddContext(ctx, handler.Logging(&handler.AuthenticateUser{}))))
    http.Handle("/mosquitto-auth-superuser", handler.CORS(handler.AddContext(ctx, handler.Logging(&handler.AuthenticateSuperUser{}))))
    http.Handle("/mosquitto-auth-acl", handler.CORS(handler.AddContext(ctx, handler.Logging(&handler.Authorize{}))))

    logger.Infof("Listening on %s...", c.GlobalString("bind-address"))
    err := http.ListenAndServe(c.GlobalString("bind-address"), nil)
    FatalOnError(err, "Failed to bind on %s: ", c.GlobalString("bind-address"))
}

func FatalOnError(err error, msg string, args ...interface{}) {
    if err != nil {
        log.Fatalf(msg, args...)
        os.Exit(1)
    }
}

func main() {
    app := cli.NewApp()

    app.Name = "PIOT Mosquitto Authentication and Authorization Server"
    app.Version = config.VersionString()
    app.Authors = []cli.Author{
        {
            Name:  "Michal Nezerka",
            Email: "michal.nezerka@gmail.com",
        },
    }
    app.Action = runServer
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:   "bind-address,b",
            Usage:  "Listen address for API HTTP endpoint",
            Value:  "0.0.0.0:9095",
            EnvVar: "BIND_ADDRESS",
        },
        cli.StringFlag{
            Name:   "mongodb-uri,m",
            Usage:  "URI for the mongo database",
            Value:  "mongodb://localhost:27017",
            EnvVar: "MONGODB_URI",
        },
        cli.StringFlag{
            Name:   "log-level,l",
            Usage:  "Logging level",
            Value:  "INFO",
            EnvVar: "LOG_LEVEL",
        },
        cli.StringFlag{
            Name:   "test-password",
            Usage:  "Test user password",
            Value:  "test",
            EnvVar: "TEST_PASSWORD",
        },
        cli.StringFlag{
            Name:   "mon-password",
            Usage:  "Monitoring user password",
            Value:  "mon",
            EnvVar: "MON_PASSWORD",
        },
        cli.StringFlag{
            Name:   "piot-password",
            Usage:  "Piot user password",
            Value:  "piot",
            EnvVar: "PIOT_PASSWORD",
        },
    }

    app.Run(os.Args)
}
