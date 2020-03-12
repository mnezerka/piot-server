package context

import (
    "context"
)

const LOG_FORMAT = "%{color}%{time:2006/01/02 15:04:05 -07:00 MST} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"

//func NewContext(dbUri, dbName string, mqtt service.IMqtt, logLevel string) context.Context {

func NewContext(o *ContextOptions) context.Context {

    // create global context for all handlers
    ctx := context.Background()

    /////////////// Parameters
    ctx = context.WithValue(ctx, "params", o.Params)

    /*

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

    /////////////// PIOT MYSQLDB SERVICE

    // create global mysql db service for all handlers
    var mysqldb service.IMysqlDb
    if o.MysqlDbHost == "mock" {
        mysqldb = &service.MysqlDbMock{}
    } else {
        // real mysqldb service instance
        mysqldb = service.NewMysqlDb(o.MysqlDbHost, o.MysqlDbUsername, o.MysqlDbPassword, o.MysqlDbName)
        err := mysqldb.Open(ctx)
        if err != nil {
            log.Fatalf("Connect to mysql server failed %v", err)
            os.Exit(1)
        }
    }
    ctx = context.WithValue(ctx, "mysqldb", mysqldb)

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

    */

    return ctx
}
