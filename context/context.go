package context

import (
    "context"
    "log"
    "os"
    "piot-server/service"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const LOG_FORMAT = "%{color}%{time:2006/01/02 15:04:05 -07:00 MST} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"

//func NewContext(dbUri, dbName string, mqtt service.IMqtt, logLevel string) context.Context {

func NewContext(o *ContextOptions) context.Context {

    // create global context for all handlers
    ctx := context.Background()

    /////////////// Parameters
    ctx = context.WithValue(ctx, "params", o.Params)

    /////////////// DB

    // try to open database
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

    /////////////// LOGGER

    // create global logger for all handlers
    log, err := service.NewLogger(LOG_FORMAT, o.Params.LogLevel)
    if err != nil {
        log.Fatalf("Cannot create logger for level %s (%v)", o.Params.LogLevel, err)
        os.Exit(1)
    }
    ctx = context.WithValue(ctx, "log", log)

    /////////////// THINGS

    // create global things service for all handlers
    things := &service.Things{}
    ctx = context.WithValue(ctx, "things", things)

    /////////////// USERS SERVICE

    // create global users service for all handlers
    users := &service.Users{}
    ctx = context.WithValue(ctx, "users", users)

    /////////////// ORGS

    // create global orgs service for all handlers
    orgs := &service.Orgs{}
    ctx = context.WithValue(ctx, "orgs", orgs)

    /////////////// AUTH

    // create global orgs service for all handlers
    auth:= &service.Auth{}
    ctx = context.WithValue(ctx, "auth", auth)

    /////////////// PIOT DEVICES SERVICE

    // create global piot devices service for all handlers
    piotdevices := service.NewPiotDevices()
    ctx = context.WithValue(ctx, "piotdevices", piotdevices)

    /////////////// PIOT INFLUXDB SERVICE

    // create global influxdb service for all handlers
    var influxdb service.IInfluxDb
    if o.InfluxDbUri == "mock" {
        influxdb = &service.InfluxDbMock{}
    } else {
        // real influxdb service instance
        influxdb = service.NewInfluxDb(o.InfluxDbUri, o.InfluxDbUsername, o.InfluxDbPassword)
    }
    ctx = context.WithValue(ctx, "influxdb", influxdb)

    /////////////// HTTP CLIENT SERVICE

    // create global http client service to be used by handlers
    var httpClient service.IHttpClient
    httpClient = service.NewHttpClient()
    ctx = context.WithValue(ctx, "httpclient", httpClient)

    /////////////// PIOT MQTT SERVICE

    // create global mqtt service for all handlers
    var mqtt service.IMqtt
    if o.MqttUri == "mock" {
        mqtt = &service.MqttMock{}
    } else {
        // mqtt instance
        mqtt = service.NewMqtt(o.MqttUri)
        mqtt.SetUsername(o.MqttUsername)
        mqtt.SetPassword(o.MqttPassword)
        mqtt.SetClient(o.MqttClient)
        err := mqtt.Connect(ctx)
        if err != nil {
            log.Fatalf("Connect to mqtt server failed %v", err)
            os.Exit(1)
        }
    }
    ctx = context.WithValue(ctx, "mqtt", mqtt)

    return ctx
}
