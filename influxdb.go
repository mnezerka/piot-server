package main

/*
https://github.com/influxdata/line-protocol/blob/master/metric.go
dlfjsdlfjkkkaaaajjjkkkttps://github.com/influxdata/influxdb-client-go/blob/develop/write.go
*/

import (
	"bytes"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	proto "github.com/influxdata/line-protocol"
	"github.com/op/go-logging"
)

type IInfluxDb interface {
	PostMeasurement(thing *Thing, value string)
	PostSwitchState(thing *Thing, value string)
	PostLocation(thing *Thing, lat, lng float64, sat, ts int32)
	PostBatteryLevel(thing *Thing, level int32)
}

type InfluxDb struct {
	log        *logging.Logger
	orgs       *Orgs
	httpClient IHttpClient
	Uri        string
	Username   string
	Password   string
}

type RowMetric struct {
	name   string
	tags   []*proto.Tag
	fields []*proto.Field
	ts     time.Time
}

func NewInfluxDb(log *logging.Logger, orgs *Orgs, httpClient IHttpClient, uri, username, password string) IInfluxDb {
	db := &InfluxDb{log: log, orgs: orgs, httpClient: httpClient}
	db.Uri = uri
	db.Username = username
	db.Password = password

	return db
}

func InfluxDbEscapeString(str string) string {

	return strings.ReplaceAll(str, " ", "\\ ")
}

/*
func paramToString(name string, value interface{}) (string, error) {
	valueStr, err := PrimitiveToString(value)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s=%s", InfluxDbEscapeString(name), InfluxDbEscapeString(valueStr))
	return result, nil
}
*/

func (db *InfluxDb) PostMeasurement(thing *Thing, value string) {
	db.log.Debugf("Posting measurement to InfluxDB, thing: %s, val: %s", thing.Name, value)

	// get thing org -> get influxdb assigned to org
	org, err := db.orgs.Get(thing.OrgId)
	if err != nil {
		return
	}

	db.log.Debugf("Going to post to InfluxDB %s as %s", org.InfluxDb, org.InfluxDbUsername)

	if thing.Type != THING_TYPE_SENSOR {
		// ignore things which don't represent sensor
		return
	}

	// get thing name, use alias if set
	name := thing.Name
	if thing.Alias != "" {
		name = thing.Alias
	}

	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		db.log.Warningf("Ignoring write measurement to influxdb for device %s due to invalid float value %s", thing.Id.Hex(), value)
		return
	}
	fields := map[string]interface{}{"value": valueFloat}
	tags := map[string]string{"id": thing.Id.Hex(), "name": name, "class": thing.Sensor.Class}
	rm := NewRowMetric("sensor", tags, fields, time.Now())
	body, err := rm.Encode()
	if err != nil {
		db.log.Errorf("Cannot encode tags and fields into InfluxDB line protocol format: %s", err.Error())
		return
	}
	//body := fmt.Sprintf("sensor,id=%s,name=%s,class=%s value=%s", thing.Id.Hex(), name, thing.Sensor.Class, value)

	url, err := url.Parse(db.Uri)
	if err != nil {
		db.log.Errorf("Cannot decode InfluxDB url from %s (%s)", db.Uri, err.Error())
		return
	}

	url.Path = path.Join(url.Path, "write")

	params := url.Query()
	params.Add("db", org.InfluxDb)
	url.RawQuery = params.Encode()

	db.httpClient.PostString(url.String(), body.String(), &db.Username, &db.Password)
}

func (db *InfluxDb) PostSwitchState(thing *Thing, value string) {
	db.log.Debugf("Posting switch state to InfluxDB, thing: %s, val: %s", thing.Name, value)

	// get thing org -> get influxdb assigned to org
	org, err := db.orgs.Get(thing.OrgId)
	if err != nil {
		return
	}

	db.log.Debugf("Going to post to InfluxDB %s as %s", org.InfluxDb, org.InfluxDbUsername)

	if thing.Type != THING_TYPE_SWITCH {
		// ignore things which don't represent switch
		return
	}

	// get thing name, use alias if set
	name := thing.Name
	if thing.Alias != "" {
		name = thing.Alias
	}

	fields := map[string]interface{}{"value": value}
	tags := map[string]string{"id": thing.Id.Hex(), "name": name}
	rm := NewRowMetric("switch", tags, fields, time.Now())
	body, err := rm.Encode()
	if err != nil {
		db.log.Errorf("Cannot encode tags and fields into InfluxDB line protocol format: %s", err.Error())
		return
	}

	//body := fmt.Sprintf("switch,id=%s,name=%s value=%s", thing.Id.Hex(), name, value)

	url, err := url.Parse(db.Uri)
	if err != nil {
		db.log.Errorf("Cannot decode InfluxDB url from %s (%s)", db.Uri, err.Error())
		return
	}

	url.Path = path.Join(url.Path, "write")

	params := url.Query()
	params.Add("db", org.InfluxDb)
	url.RawQuery = params.Encode()

	db.httpClient.PostString(url.String(), body.String(), &db.Username, &db.Password)
}

func (db *InfluxDb) PostLocation(thing *Thing, lat, lng float64, sat, ts int32) {
	db.log.Debugf("Posting thing location to InfluxDB, thing: %s, lat: %f, lng: %f, sat: %d, ts: %d", thing.Name, lat, lng, sat, ts)

	// get thing org -> get influxdb assigned to org
	org, err := db.orgs.Get(thing.OrgId)
	if err != nil {
		return
	}

	db.log.Debugf("Going to post location to InfluxDB %s as %s", org.InfluxDb, org.InfluxDbUsername)

	// get thing name, use alias if set
	name := thing.Name
	if thing.Alias != "" {
		name = thing.Alias
	}

	fields := map[string]interface{}{
		"lat": lat,
		"lng": lng,
		"sat": int64(sat),
	}
	tags := map[string]string{
		"id":   thing.Id.Hex(),
		"name": name,
	}

	rm := NewRowMetric("location", tags, fields, time.Unix(int64(ts), 0))
	buf, err := rm.Encode()
	if err != nil {
		db.log.Errorf("Cannot encode tags and fields into InfluxDB line protocol format: %s", err.Error())
		return
	}

	url, err := url.Parse(db.Uri)
	if err != nil {
		db.log.Errorf("Cannot decode InfluxDB url from %s (%s)", db.Uri, err.Error())
		return
	}

	url.Path = path.Join(url.Path, "write")

	params := url.Query()
	params.Add("db", org.InfluxDb)
	url.RawQuery = params.Encode()

	db.httpClient.PostString(url.String(), buf.String(), &db.Username, &db.Password)
}

func (db *InfluxDb) PostBatteryLevel(thing *Thing, level int32) {
	db.log.Debugf("Posting thing battery level to InfluxDB, thing: %s, level: %d", thing.Name, level)

	// get thing org -> get influxdb assigned to org
	org, err := db.orgs.Get(thing.OrgId)
	if err != nil {
		return
	}

	db.log.Debugf("Going to post to InfluxDB %s as %s", org.InfluxDb, org.InfluxDbUsername)

	// get thing name, use alias if set
	name := thing.Name
	if thing.Alias != "" {
		name = thing.Alias
	}

	fields := map[string]interface{}{"level": int64(level)}
	tags := map[string]string{"id": thing.Id.Hex(), "name": name}
	rm := NewRowMetric("battery", tags, fields, time.Now())
	body, err := rm.Encode()
	if err != nil {
		db.log.Errorf("Cannot encode tags and fields into InfluxDB line protocol format: %s", err.Error())
		return
	}

	url, err := url.Parse(db.Uri)
	if err != nil {
		db.log.Errorf("Cannot decode InfluxDB url from %s (%s)", db.Uri, err.Error())
		return
	}

	url.Path = path.Join(url.Path, "write")

	params := url.Query()
	params.Add("db", org.InfluxDb)
	url.RawQuery = params.Encode()

	db.httpClient.PostString(url.String(), body.String(), &db.Username, &db.Password)
}

func NewRowMetric(
	name string,
	tags map[string]string,
	fields map[string]interface{},
	ts time.Time,
) *RowMetric {
	m := &RowMetric{
		name:   name,
		tags:   nil,
		fields: nil,
		ts:     ts,
	}

	// convert tags to protocol format
	if len(tags) > 0 {
		m.tags = make([]*proto.Tag, 0, len(tags))
		for k, v := range tags {
			m.tags = append(m.tags,
				&proto.Tag{Key: k, Value: v})
		}
	}

	// convert fields to protocol format
	m.fields = make([]*proto.Field, 0, len(fields))
	for k, v := range fields {
		/*v := convertField(v)
		  if v == nil {
		      continue
		  }
		*/
		m.fields = append(m.fields, &proto.Field{Key: k, Value: v})
	}
	return m
}

func (rm *RowMetric) Time() time.Time           { return rm.ts }
func (rm *RowMetric) Name() string              { return rm.name }
func (rm *RowMetric) TagList() []*proto.Tag     { return rm.tags }
func (rm *RowMetric) FieldList() []*proto.Field { return rm.fields }

func (rm *RowMetric) Encode() (*bytes.Buffer, error) {

	buf := &bytes.Buffer{}
	e := proto.NewEncoder(buf)
	e.SetFieldTypeSupport(proto.UintSupport)
	e.SetFieldSortOrder(proto.SortFields)
	e.FailOnFieldErr(true)

	if _, err := e.Encode(rm); err != nil {
		return nil, err
	}

	return buf, nil
}
