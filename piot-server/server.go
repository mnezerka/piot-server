/*
 Resources:
    https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/
*/

package main

import (
    "log"
    "net/http"
    "os"
    "github.com/urfave/cli"
    "piot-server/handler"
    "piot-server/config"
    "piot-server/resolver"
    "piot-server/schema"
    "piot-server/test"
    piotcontext "piot-server/context"
    graphql "github.com/graph-gophers/graphql-go"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/op/go-logging"
)

func runServer(c *cli.Context) {

    // create global context for all handlers
    ctx := piotcontext.NewContext(c.GlobalString("mongodb-uri"), "piot", &test.MqttMock{}, c.GlobalString("log-level"))

    // Auto disconnect from mongo
    defer ctx.Value("dbClient").(*mongo.Client).Disconnect(ctx)

    logger := ctx.Value("log").(*logging.Logger)
    logger.Infof("Starting PIOT server %s", config.VersionString())

    /////////////// HANDLERS

    // create GraphQL schema
    graphqlSchema := graphql.MustParseSchema(schema.GetRootSchema(), &resolver.Resolver{})

    http.HandleFunc("/", handler.RootHandler)

    // endpoint for registration of new user
    http.Handle("/register", handler.CORS(handler.AddContext(ctx, handler.Logging(handler.Registration()))))

    // endpoint for authentication - token is generaged
    http.Handle("/login", handler.CORS(handler.AddContext(ctx, handler.Logging(handler.LoginHandler()))))

    // endpoint for refreshing nearly expired token
    //r.HandleFunc("/refresh", handler.RefreshHandler)

    http.Handle("/query", handler.CORS(handler.AddContext(ctx, handler.Logging(handler.Authorize(&handler.GraphQL{Schema: graphqlSchema})))))
    //http.Handle("/query", handler.AddContext(ctx, handler.Logging(&handler.GraphQL{Schema: graphqlSchema})))

    // enpoint for interactive graphql web IDE
    http.Handle("/gql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "graphiql.html")
    }))

    //http.Handle("/adapter", handler.CORS(handler.AddContext(ctx, handler.Logging(handler.Authorize(&handler.Adapter{})))))
    http.Handle("/adapter", handler.CORS(handler.AddContext(ctx, handler.Logging(&handler.Adapter{}))))

    logger.Infof("Listening on %s...", c.GlobalString("bind-address"))
    //err = http.ListenAndServe(c.GlobalString("bind-address"), handlers.LoggingHandler(os.Stdout, r))
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
        cli.StringFlag{
            Name:   "log-level,l",
            Usage:  "Logging level",
            Value:  "INFO",
            EnvVar: "LOG_LEVEL",
        },
    }

    app.Run(os.Args)
}
