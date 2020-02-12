package service

import (
    "context"
    "piot-server/model"
    "github.com/op/go-logging"
)

type IMysqlDb interface {
    StoreMeasurement(ctx context.Context, thing *model.Thing, value string)
    StoreSwitchState(ctx context.Context, thing *model.Thing, value string)
}

type MysqlDb struct {
    Host string
    Username string
    Password string
}

func NewMysqlDb(host, username, password string) IMysqlDb {
    db := &MysqlDb{}
    db.Host = host
    db.Username = username
    db.Password = password

    return db
}

func (db *MysqlDb) StoreMeasurement(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Storing measurement to mysql db, thing: %s, val: %s", thing.Name, value)

    // get thing org -> get mysql db assigned to org
    orgs := ctx.Value("orgs").(*Orgs)
    org, err := orgs.Get(ctx, thing.OrgId)
    if err != nil {
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Going to post to mysql db %s as %s", org.MysqlDb, org.MysqlDbUsername)

    // TODO
}

func (db *MysqlDb) StoreSwitchState(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Storing switch state to MysqlDb, thing: %s, val: %s", thing.Name, value)

    // get thing org -> get influxdb assigned to org
    orgs := ctx.Value("orgs").(*Orgs)
    org, err := orgs.Get(ctx, thing.OrgId)
    if err != nil {
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Going to post to InfluxDB %s as %s", org.MysqlDb, org.MysqlDbUsername)

    // TODO
}
