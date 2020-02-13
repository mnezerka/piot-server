package service

import (
    "context"
    "fmt"
    "strconv"
    "time"
    "piot-server/model"
    "github.com/op/go-logging"
    "database/sql"
    _"github.com/go-sql-driver/mysql"
)

type IMysqlDb interface {
    Open(ctx context.Context) error
    Close(ctx context.Context)
    StoreMeasurement(ctx context.Context, thing *model.Thing, value string)
    StoreSwitchState(ctx context.Context, thing *model.Thing, value string)
}

type MysqlDb struct {
    Host string
    Username string
    Password string
    Db *sql.DB
}

func NewMysqlDb(host, username, password string) IMysqlDb {
    db := &MysqlDb{}
    db.Host = host
    db.Username = username
    db.Password = password
    db.Db = nil

    return db
}

func (db *MysqlDb) Open(ctx context.Context) error {
    ctx.Value("log").(*logging.Logger).Infof("Connecting to mysql database %s", db.Host)

    // open database if host is specified
    if db.Host == "" {
        ctx.Value("log").(*logging.Logger).Warningf("Refusing to open mysql database, host not specified")
        return nil
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", db.Username, db.Password, db.Host, "piot")
    ctx.Value("log").(*logging.Logger).Debugf("Mysql DSN: %s", dsn)

    d, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    // Open doesn't open a connection. Validate DSN data:
    err = d.Ping()
    if err != nil {
        return err
    }
    db.Db = d

    ctx.Value("log").(*logging.Logger).Infof("Connected to mysql database")

    return nil
}

func (db *MysqlDb) Close(ctx context.Context) {
    if db.Db != nil {
        db.Db.Close()
    }
}

func (db *MysqlDb) verifyOrg(ctx context.Context, thing *model.Thing) *model.Org {

    if db.Db == nil {
        ctx.Value("log").(*logging.Logger).Warningf("Mysql database is not initialized")
        return nil
    }

    // get thing org
    orgs := ctx.Value("orgs").(*Orgs)
    org, err := orgs.Get(ctx, thing.OrgId)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Warningf("Store to database rejected for thing (%s) that is not assigned to org", thing.Id.Hex())
        return nil
    }

    // mysql name needs to be configured
    if org.MysqlDb == "" {
        ctx.Value("log").(*logging.Logger).Warningf("Store to database rejected for thing (%s), where org (%s) has no mysql configuration", thing.Id.Hex(), org.Name)
        return nil
    }

    return org
}


func (db *MysqlDb) StoreMeasurement(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Storing measurement to mysql db, thing: %s, val: %s", thing.Name, value)

    org := db.verifyOrg(ctx, thing)

    // convert value to float
    valueFloat, err := strconv.ParseFloat(value, 32)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Mysql database storage - float conversion error for value %s", value)
        return
    }

    query := "INSERT IGNORE INTO piot_sensors (`id`, `org`, `class`, `value`, `time`) VALUES (?, ?, ?, ?, ?)"

    r, err := db.Db.Query(query, thing.Id.Hex(), org.MysqlDb, thing.Sensor.Class, valueFloat, int32(time.Now().Unix()))

    // Failure when trying to store data
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Mysql database operation failed: %s", err.Error())
    }

    r.Close() // Always do this or you will leak connections
}

func (db *MysqlDb) StoreSwitchState(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Storing switch state to MysqlDb, thing: %s, val: %s", thing.Name, value)

    org := db.verifyOrg(ctx, thing)

    if thing.Type != model.THING_TYPE_SWITCH {
        // ignore things which don't represent switch
        return
    }

    // convert value to int 
    valueInt, err := strconv.Atoi(value)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Mysql database storage - int conversion error for value %s", value)
        return
    }

    query := "INSERT IGNORE INTO piot_switches (`id`, `org`, `value`, `time`) VALUES (?, ?, ?, ?)"

    r, err := db.Db.Query(query, thing.Id.Hex(), org.MysqlDb, valueInt, int32(time.Now().Unix()))

    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Mysql database operation failed: %s", err.Error())
    }

    r.Close() // Always do this or you will leak connections
}
