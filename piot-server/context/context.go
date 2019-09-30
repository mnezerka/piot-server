package context

import (
    "log"
    "os"
    "golang.org/x/net/context"
    "piot-server/service"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const LOG_FORMAT = "%{color}%{time:2006/01/02 15:04:05 -07:00 MST} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"

func NewContext(dbUri string, dbName string) context.Context {

    // create global context for all handlers
    ctx := context.Background()

    /////////////// DB

    // try to open database
    dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
    if err != nil {
        log.Fatalf("Failed to open database on %s (%v)", dbUri, err)
        os.Exit(1)
    }

    // Check the connection
    err = dbClient.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("Cannot ping database on %s (%v)", dbUri, err)
        os.Exit(1)
    }

    // Auto disconnect from mongo
    //defer dbClient.Disconnect(ctx)

    ctx = context.WithValue(ctx, "dbClient", dbClient)

    db := dbClient.Database(dbName)
    ctx = context.WithValue(ctx, "db", db)

    /////////////// LOGGER

    // create global logger for all handlers
    log := service.NewLogger(LOG_FORMAT, true)
    ctx = context.WithValue(ctx, "log", log)

    return ctx
}