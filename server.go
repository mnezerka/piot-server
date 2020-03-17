package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"
    "github.com/urfave/cli"
    "piot-server/handler"
    piotconfig "piot-server/config"
    "piot-server/resolver"
    "piot-server/schema"
    "github.com/mnezerka/go-piot"
    "github.com/mnezerka/go-piot/config"
    piotcontext "piot-server/context"
    graphql "github.com/graph-gophers/graphql-go"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    //"github.com/op/go-logging"
)

const LOG_FORMAT = "%{color}%{time:2006/01/02 15:04:05 -07:00 MST} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"

func runServer(c *cli.Context) {

    cfg := config.NewParameters()
    cfg.DbUri = c.GlobalString("mongodb-uri")
    cfg.DbName = "piot"
    cfg.LogLevel = c.GlobalString("log-level")

    ///////////////// LOGGER instance
    logger, err := piot.NewLogger(LOG_FORMAT, cfg.LogLevel)
    if err != nil {
        log.Fatalf("Cannot create logger for level %s (%v)", cfg.LogLevel, err)
        os.Exit(1)
    }

    /////////////// DB (mongo)
    dbUri := c.GlobalString("mongodb-uri")

    // try to open database
    dbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
    if err != nil {
        logger.Fatalf("Failed to open database on %s (%v)", dbUri, err)
        os.Exit(1)
    }

    // Check the connection
    err = dbClient.Ping(context.TODO(), nil)
    if err != nil {
        logger.Fatalf("Cannot ping database on %s (%v)", dbUri, err)
        os.Exit(1)
    }

    db := dbClient.Database("piot")

    /////////////// USERS service
    users := piot.NewUsers(logger, db)

    /////////////// ORGS service
    orgs := piot.NewOrgs(logger, db)

    /////////////// HTTP CLIENT service
    var httpClient piot.IHttpClient
    httpClient = piot.NewHttpClient(logger)

    /////////////// PIOT INFLUXDB SERVICE

    // create global influxdb service for all handlers
    influxDbUri := c.GlobalString("influxdb-uri")
    influxDbUsername := c.GlobalString("influxdb-user")
    influxDbPassword := c.GlobalString("influxdb-password")
    influxDb := piot.NewInfluxDb(logger, orgs, httpClient, influxDbUri, influxDbUsername, influxDbPassword)

    /////////////// PIOT MYSQLDB SERVICE

    // create global mysql db service for all handlers
    mysqlDbHost := c.GlobalString("mysqldb-host")
    mysqlDbUsername := c.GlobalString("mysqldb-user")
    mysqlDbPassword := c.GlobalString("mysqldb-password")
    mysqlDbName := c.GlobalString("mysqldb-name")

    // real mysqldb service instance
    mysqlDb := piot.NewMysqlDb(logger, orgs, mysqlDbHost, mysqlDbUsername, mysqlDbPassword, mysqlDbName)
    err = mysqlDb.Open()
    if err != nil {
        logger.Fatalf("Connect to mysql server failed %v", err)
        os.Exit(1)
    }

    //////////////// THINGS service instance
    things := piot.NewThings(db, logger)

    /////////////// PIOT MQTT service instance
    mqttUri := c.GlobalString("mqtt-uri")
    mqttUsername := c.GlobalString("mqtt-user")
    mqttPassword := c.GlobalString("mqtt-password")
    mqttClient := c.GlobalString("mqtt-client")
    mqtt := piot.NewMqtt(mqttUri, logger, things, orgs, influxDb, mysqlDb)
    mqtt.SetUsername(mqttUsername)
    mqtt.SetPassword(mqttPassword)
    mqtt.SetClient(mqttClient)
    err = mqtt.Connect(true)
    if err != nil {
        logger.Fatalf("Connect to mqtt server failed %v", err)
        os.Exit(1)
    }

    /////////////// PIOT DEVICES service instance
    piotDevices := piot.NewPiotDevices(logger, things, mqtt, cfg)

    // create global context for all handlers
    contextOptions := piotcontext.NewContextOptions()
    contextOptions.DbUri = c.GlobalString("mongodb-uri")
    contextOptions.DbName = "piot"
    contextOptions.MqttUri = c.GlobalString("mqtt-uri")
    contextOptions.MqttUsername = c.GlobalString("mqtt-user")
    contextOptions.MqttPassword = c.GlobalString("mqtt-password")
    contextOptions.MqttClient = c.GlobalString("mqtt-client")
    contextOptions.InfluxDbUri = c.GlobalString("influxdb-uri")
    contextOptions.InfluxDbUsername = c.GlobalString("influxdb-user")
    contextOptions.InfluxDbPassword = c.GlobalString("influxdb-password")
    contextOptions.MysqlDbHost = c.GlobalString("mysqldb-host")
    contextOptions.MysqlDbUsername = c.GlobalString("mysqldb-user")
    contextOptions.MysqlDbPassword = c.GlobalString("mysqldb-password")
    contextOptions.MysqlDbName = c.GlobalString("mysqldb-name")
    contextOptions.Params.LogLevel = c.GlobalString("log-level")
    contextOptions.Params.DOSInterval = c.GlobalDuration("dos-interval")
    contextOptions.Params.JwtPassword = c.GlobalString("jwt-password")
    contextOptions.Params.JwtTokenExpiration = c.GlobalDuration("jwt-token-expiration")

    ctx := piotcontext.NewContext(contextOptions)

    // Auto disconnect from mongo
    //defer ctx.Value("dbClient").(*mongo.Client).Disconnect(ctx)

    logger.Infof("Starting PIOT server %s", piotconfig.VersionString())

    /////////////// HANDLERS

    // create GraphQL schema together with resolver
    resolver := resolver.NewResolver(logger, db, orgs, users, things)
    graphqlSchema := graphql.MustParseSchema(schema.GetRootSchema(), resolver)

    http.HandleFunc("/", handler.RootHandler)

    // endpoint for registration of new user
    http.Handle("/register", handler.CORS(handler.AddContext(ctx, handler.Logging(handler.Registration()))))

    // endpoint for authentication - token is generaged
    loginHandler := handler.NewLogin(logger, db, cfg)
    http.Handle("/login", handler.CORS(handler.AddContext(ctx, handler.Logging(loginHandler))))

    // endpoint for refreshing nearly expired token
    //r.HandleFunc("/refresh", handler.RefreshHandler)

    http.Handle("/query", handler.CORS(handler.AddContext(ctx, handler.Logging(handler.Authorize(&handler.GraphQL{Schema: graphqlSchema})))))
    //http.Handle("/query", handler.AddContext(ctx, handler.Logging(&handler.GraphQL{Schema: graphqlSchema})))

    // enpoint for interactive graphql web IDE
    http.Handle("/gql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "graphiql.html")
    }))

    //http.Handle("/adapter", handler.CORS(handler.AddContext(ctx, handler.Logging(handler.Authorize(&handler.Adapter{})))))
    adapterHandler := handler.NewAdapter(logger, piotDevices)
    http.Handle("/adapter", handler.CORS(handler.AddContext(ctx, handler.Logging(adapterHandler))))

    logger.Infof("Listening on %s...", c.GlobalString("bind-address"))
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
    app.Version = piotconfig.VersionString()
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
        cli.StringFlag{
            Name:   "mqtt-user",
            Usage:  "Username for mqtt authentication",
            EnvVar: "MQTT_USER",
        },
        cli.StringFlag{
            Name:   "mqtt-password",
            Usage:  "Password for mqtt authentication",
            EnvVar: "MQTT_PASSWORD",
        },
        cli.StringFlag{
            Name:   "mqtt-client",
            Usage:  "Id used for identification of this mqtt client",
            Value:  "piot-server",
            EnvVar: "MQTT_CLIENT",
        },
        cli.DurationFlag{
            Name: "dos-interval",
            Usage: "The minimal allowed time interval between two packets from the same device",
            Value: time.Second * 1,
            EnvVar: "DOS_INTERVAL",
        },
        cli.StringFlag{
            Name:   "jwt-password",
            Usage:  "Password for jwt communication",
            EnvVar: "JWT_PASSWORD",
            Value: "secret-key",
        },
        cli.DurationFlag{
            Name: "jwt-token-expiration",
            Usage: "Expriation of JWT token in seconds",
            Value: time.Hour * 4,
            EnvVar: "JWT_TOKEN_EXPIRATION",
        },
        cli.StringFlag{
            Name:   "influxdb-uri",
            Usage:  "URI for the InfluxDB database",
            EnvVar: "INFLUXDB_URI",
        },
        cli.StringFlag{
            Name:   "influxdb-user",
            Usage:  "Username for InfluxDB user with admin privileges",
            EnvVar: "INFLUXDB_USER",
        },
        cli.StringFlag{
            Name:   "influxdb-password",
            Usage:  "Password for InfluxDB user with admin privileges",
            EnvVar: "INFLUXDB_PASSWORD",
        },
        cli.StringFlag{
            Name:   "mysqldb-host",
            Usage:  "Hostname for the Mysql database",
            EnvVar: "MYSQLDB_HOST",
        },
        cli.StringFlag{
            Name:   "mysqldb-user",
            Usage:  "Username for mysql user with admin privileges",
            EnvVar: "MYSQLDB_USER",
        },
        cli.StringFlag{
            Name:   "mysqldb-password",
            Usage:  "Password for mysql user with admin privileges",
            EnvVar: "MYSQLDB_PASSWORD",
        },
        cli.StringFlag{
            Name:   "mysqldb-name",
            Usage:  "Mysql database name",
            EnvVar: "MYSQLDB_NAME",
        },
    }

    app.Run(os.Args)
}
