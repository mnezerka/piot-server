package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
)

type IMysqlDb interface {
	Open() error
	Close()
	StoreMeasurement(thing *Thing, value string)
	StoreSwitchState(thing *Thing, value string)
}

type MysqlDb struct {
	log      *logging.Logger
	orgs     *Orgs
	Host     string
	Username string
	Password string
	Name     string
	Db       *sql.DB
}

func NewMysqlDb(log *logging.Logger, orgs *Orgs, host, username, password, name string) IMysqlDb {
	db := &MysqlDb{log: log, orgs: orgs}
	db.Host = host
	db.Username = username
	db.Password = password
	db.Name = name
	db.Db = nil

	return db
}

func (db *MysqlDb) Open() error {
	db.log.Infof("Connecting to mysql database %s", db.Host)

	// open database if host is specified
	if db.Host == "" || db.Name == "" {
		db.log.Warningf("Refusing to open mysql database, host or db name not specified")
		return nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", db.Username, db.Password, db.Host, db.Name)
	db.log.Debugf("Mysql DSN: %s", dsn)

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

	db.log.Infof("Connected to mysql database")

	return nil
}

func (db *MysqlDb) Close() {
	if db.Db != nil {
		db.Db.Close()
	}
}

func (db *MysqlDb) verifyOrg(thing *Thing) *Org {
	if db.Db == nil {
		db.log.Warningf("Mysql database is not initialized")
		return nil
	}

	// get thing org
	org, err := db.orgs.Get(thing.OrgId)
	if err != nil {
		db.log.Warningf("Store to database rejected for thing (%s) that is not assigned to org", thing.Id.Hex())
		return nil
	}

	// mysql name needs to be configured
	if org.MysqlDb == "" {
		db.log.Warningf("Store to database rejected for thing (%s), where org (%s) has no mysql configuration", thing.Id.Hex(), org.Name)
		return nil
	}

	return org
}

func (db *MysqlDb) getTimestamp(thing *Thing) int32 {
	// generate unix timestamp
	ts := int32(time.Now().Unix())

	// alter timestamp to match low boundary of configured interval
	if thing.StoreMysqlDbInterval > 0 {
		ts = ts - (ts % thing.StoreMysqlDbInterval)
	}

	return ts
}

func (db *MysqlDb) StoreMeasurement(thing *Thing, value string) {
	db.log.Debugf("Storing measurement to mysql db, thing: %s, val: %s", thing.Name, value)

	// verify if all preconditions are met
	org := db.verifyOrg(thing)
	if org == nil {
		return
	}

	// convert value to float
	valueFloat, err := strconv.ParseFloat(value, 32)
	if err != nil {
		db.log.Errorf("Mysql database storage - float conversion error for value %s", value)
		return
	}

	ts := db.getTimestamp(thing)

	query := "INSERT IGNORE INTO piot_sensors (`id`, `org`, `class`, `value`, `time`) VALUES (?, ?, ?, ?, ?)"

	r, err := db.Db.Query(query, thing.Id.Hex(), org.MysqlDb, thing.Sensor.Class, valueFloat, ts)

	// Failure when trying to store data
	if err != nil {
		db.log.Errorf("Mysql database operation failed: %s", err.Error())
	}

	r.Close() // Always do this or you will leak connections
}

func (db *MysqlDb) StoreSwitchState(thing *Thing, value string) {
	db.log.Debugf("Storing switch state to MysqlDb, thing: %s, val: %s", thing.Name, value)

	// verify if all preconditions are met
	org := db.verifyOrg(thing)
	if org == nil {
		return
	}

	if thing.Type != THING_TYPE_SWITCH {
		// ignore things which don't represent switch
		return
	}

	// convert value to int
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		db.log.Errorf("Mysql database storage - int conversion error for value %s", value)
		return
	}

	ts := db.getTimestamp(thing)

	query := "INSERT IGNORE INTO piot_switches (`id`, `org`, `value`, `time`) VALUES (?, ?, ?, ?)"

	r, err := db.Db.Query(query, thing.Id.Hex(), org.MysqlDb, valueInt, ts)

	if err != nil {
		db.log.Errorf("Mysql database operation failed: %s", err.Error())
	}

	r.Close() // Always do this or you will leak connections
}
