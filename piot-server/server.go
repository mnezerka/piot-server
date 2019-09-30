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
    "piot-server/handler"
    "piot-server/config"
    "piot-server/service"
    "piot-server/resolver"
    "piot-server/schema"
    graphql "github.com/graph-gophers/graphql-go"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const LOG_FORMAT = "%{color}%{time:2006/01/02 15:04:05 -07:00 MST} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"

func runServer(c *cli.Context) {

    // create global context for all handlers
    ctx := context.Background()

    /////////////// DB
    // try to open database
    dbUri := c.GlobalString("mongodb-uri")
    dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
    FatalOnError(err, "Failed to open database on %s", dbUri)

    // Check the connection
    err = dbClient.Ping(ctx, nil)
    FatalOnError(err, "Cannot ping database on %s", dbUri)

    // Auto disconnect from mongo
    defer dbClient.Disconnect(ctx)

    db := dbClient.Database("piot")
    ctx = context.WithValue(ctx, "db", db)

    /////////////// LOGGER

    // create global logger for all handlers
    log := service.NewLogger(LOG_FORMAT, true)
    ctx = context.WithValue(ctx, "log", log)

    log.Infof("Starting PIOT server %s", config.VersionString())

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

    http.Handle("/adapter", handler.CORS(handler.AddContext(ctx, handler.Logging(handler.Authorize(&handler.Adapter{})))))

    log.Infof("Listening on %s...", c.GlobalString("bind-address"))
    //err = http.ListenAndServe(c.GlobalString("bind-address"), handlers.LoggingHandler(os.Stdout, r))
    err = http.ListenAndServe(c.GlobalString("bind-address"), nil)
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
    }

    app.Run(os.Args)
}
