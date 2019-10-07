package context

import (
    "context"
    "os"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "mosquitto-auth/service"
)

const LOG_FORMAT = "%{color}%{time:2006/01/02 15:04:05 -07:00 MST} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"

func NewContext(o *ContextOptions) context.Context {

    // create global context for all handlers
    ctx := context.Background()

    /////////////// LOGGER

    // create global logger for all handlers
    log, err := service.NewLogger(LOG_FORMAT, o.LogLevel)
    if err != nil {
        log.Fatalf("Cannot create logger for level %s (%v)", o.LogLevel, err)
        os.Exit(1)
    }
    ctx = context.WithValue(ctx, "log", log)

    /////////////// STATIC USERS

    ctx = context.WithValue(ctx, "test-pwd", o.TestPassword)
    ctx = context.WithValue(ctx, "mon-pwd", o.MonPassword)
    ctx = context.WithValue(ctx, "piot-pwd", o.PiotPassword)

    /////////////// DB

    // try to open database
    log.Debugf("Connecting to mongodb database <%s>", o.DbUri)
    dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(o.DbUri))
    if err != nil {
        log.Fatalf("Failed to open database on %s (%v)", o.DbUri, err)
        os.Exit(1)
    }

    // Check the connection
    err = dbClient.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("Cannot ping database on %s (%v)", o.DbUri, err)
        os.Exit(1)
    }

    ctx = context.WithValue(ctx, "dbClient", dbClient)

    db := dbClient.Database(o.DbName)
    ctx = context.WithValue(ctx, "db", db)

    log.Debugf("Connected to mongodb database <%s>", o.DbUri)

    return ctx
}


